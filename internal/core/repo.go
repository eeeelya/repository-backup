package core

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/term"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Repo struct {
	URL string
}

func (r *Repo) GetRepoName() string {
	parts := strings.Split(strings.TrimSuffix(r.URL, ".git"), "/")

	return parts[len(parts)-1]
}

func getPassword() (string, error) {
	fmt.Print("Enter SSH key password: ")

	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Println()

	return string(passwordBytes), nil
}

func (r *Repo) getAuthMethod(sshKeyPath string) (transport.AuthMethod, error) {
	if strings.HasPrefix(r.URL, "git@") || strings.HasPrefix(r.URL, "ssh://") {
		slog.Debug("Using SSH authentication")

		password, err := getPassword()
		if err != nil {
			return nil, fmt.Errorf("failed to get password: %w", err)
		}

		auth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, password)

		if err == nil {
			return auth, nil
		} else {
			return nil, fmt.Errorf("failed to load SSH keys: %w", err)
		}
	}

	slog.Debug("Using HTTPS (no authentication)")

	return nil, nil
}

type fileJob struct {
	path    string
	content []byte
}

func (r *Repo) createZipArchive(fs billy.Filesystem, outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return r.createZipArchiveConcurrent(fs, zipWriter)
}

func (r *Repo) createZipArchiveConcurrent(fs billy.Filesystem, zipWriter *zip.Writer) error {
	filePaths, err := r.CollectFilePaths(fs, "/")
	if err != nil {
		return err
	}

	if len(filePaths) == 0 {
		return nil
	}

	slog.Debug("Processing files", "count", len(filePaths), "workers", runtime.NumCPU())

	jobs := make(chan string, len(filePaths))
	results := make(chan fileJob, len(filePaths))
	errors := make(chan error, 1)

	numWorkers := min(runtime.NumCPU(), 8)

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				file, err := fs.Open(path)
				if err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}
				defer file.Close()

				var buf bytes.Buffer
				if _, err := io.Copy(&buf, file); err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}

				results <- fileJob{
					path:    path,
					content: buf.Bytes(),
				}
			}
		}()
	}

	go func() {
		for _, path := range filePaths {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		select {
		case err := <-errors:
			return err
		default:
		}

		zipEntry, err := zipWriter.Create(strings.TrimPrefix(result.path, "/"))
		if err != nil {
			return err
		}

		if _, err := io.Copy(zipEntry, bytes.NewReader(result.content)); err != nil {
			return err
		}
	}

	select {
	case err := <-errors:
		return err
	default:
	}

	return nil
}

func (r *Repo) CollectFilePaths(fs billy.Filesystem, path string) ([]string, error) {
	var paths []string

	files, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())

		if file.Name() == ".git" {
			continue
		}

		if file.IsDir() {
			subPaths, err := r.CollectFilePaths(fs, fullPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else {
			paths = append(paths, fullPath)
		}
	}

	return paths, nil
}

func (r *Repo) Save(outputPath string, sshKeyPath string) error {
	memoryStore := memory.NewStorage()
	memoryFilesystem := memfs.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	cleanSshKeyPath := strings.Replace(sshKeyPath, "~", homeDir, 1)

	slog.Debug("Cloning repository", "url", r.URL)
	slog.Debug("Using SSH key", "path", cleanSshKeyPath)

	auth, err := r.getAuthMethod(cleanSshKeyPath)
	if err != nil {
		return fmt.Errorf("failed to setup authentication: %w", err)
	}

	_, err = git.Clone(memoryStore, memoryFilesystem, &git.CloneOptions{
		URL:      r.URL,
		Progress: os.Stdout,
		Auth:     auth,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	slog.Debug("Repository cloned successfully, creating zip archive...")

	var fullOutputPath string

	cleanedPath := strings.Replace(outputPath, "~", homeDir, 1)

	absPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	fullOutputPath = filepath.Join(absPath, r.GetRepoName()+".zip")

	if err := r.createZipArchive(memoryFilesystem, fullOutputPath); err != nil {
		os.Remove(fullOutputPath)
		return fmt.Errorf("failed to create zip archive: %w", err)
	}

	slog.Debug("Backup completed successfully", "output", fullOutputPath)
	return nil
}
