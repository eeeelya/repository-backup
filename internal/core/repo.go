package core

import (
	// "archive/zip"
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"

	// "io"
	"log/slog"
	"os"
	"strings"

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

func (r *Repo) getAuthMethod(sshKeyPath string) (transport.AuthMethod, error) {
	if strings.HasPrefix(r.URL, "git@") || strings.HasPrefix(r.URL, "ssh://") {
		slog.Debug("Using SSH authentication")

		var password string
		fmt.Print("Enter ssh key password: ")
		fmt.Scan(&password)

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

func (r *Repo) createZipArchive(fs billy.Filesystem, outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return r.walkFilesystem(fs, "/", zipWriter)
}

func (r *Repo) walkFilesystem(fs billy.Filesystem, path string, zipWriter *zip.Writer) error {
	files, err := fs.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())

		if file.Name() == ".git" {
			continue
		}

		if file.IsDir() {
			if err := r.walkFilesystem(fs, fullPath, zipWriter); err != nil {
				return err
			}
		} else {
			zipEntry, err := zipWriter.Create(strings.TrimPrefix(fullPath, "/"))
			if err != nil {
				return err
			}

			fileContent, err := fs.Open(fullPath)
			if err != nil {
				return err
			}
			defer fileContent.Close()

			if _, err := io.Copy(zipEntry, fileContent); err != nil {
				return err
			}
		}
	}

	return nil
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
