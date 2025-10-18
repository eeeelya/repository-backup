# Repository Backup

A command-line tool for backing up Git repositories as zip archives. Clone any repository (public or private) and save it to a zip file with a single command.

## Installation

### Quick Install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/eeeelya/repository-backup/master/install.sh | bash
```

### Manual Installation

#### Download Pre-built Binary

1. Go to [Releases](https://github.com/eeeelya/repository-backup/releases/latest)
2. Download the binary for your platform:
   - **Linux AMD64**: `git-repository-linux-amd64`
   - **Linux ARM64**: `git-repository-linux-arm64`
   - **macOS AMD64**: `git-repository-darwin-amd64`
   - **macOS ARM64** (M1/M2): `git-repository-darwin-arm64`
   - **Windows**: `git-repository-windows-amd64.exe`

### Install from Source

Requires Go 1.21 or higher:

```bash
go install github.com/eeeelya/repository-backup/cmd/git-repository@latest
```

Or build manually:

```bash
git clone https://github.com/eeeelya/repository-backup.git
cd repository-backup
go build -o git-repository ./cmd/git-repository
sudo mv git-repository /usr/local/bin/
```

## Usage

### Basic Examples

```bash
# Save repository to current directory (creates <repo-name>.zip)
git-repository save https://github.com/user/repo .

# Save to home directory
git-repository save https://github.com/user/repo ~

# Save with custom path
git-repository save https://github.com/user/repo /tmp/backups

# Clone private repository with SSH
git-repository save git@github.com:user/private-repo.git .

# Use custom SSH key
git-repository save git@github.com:user/repo.git . --ssh-key ~/.ssh/custom_key
```

### Available Commands

```bash
git-repository version           # Show version
git-repository save <url> <path> # Backup repository
git-repository help              # Show help
git-repository help save         # Show detailed help for save command
```

### Command Options

```
git-repository save [repository-url] [output-path] [flags]

Flags:
  --ssh-key string   Path to SSH private key (default: ~/.ssh/id_rsa)
  -h, --help         Help for save
```

## SSH Authentication

For private repositories, the tool automatically looks for SSH keys in:
- `~/.ssh/id_rsa`
- `~/.ssh/id_ed25519`
- `~/.ssh/id_ecdsa`

Or specify a custom key:
```bash
git-repository save git@github.com:user/repo.git . --ssh-key ~/.ssh/my_key
```

## Uninstallation

### If Installed via Script or Manual Download

```bash
# Remove the binary
sudo rm /usr/local/bin/git-repository

# Verify removal
which git-repository  # Should return nothing
```

### If Installed via `go install`

```bash
# Remove from GOPATH/bin
rm $(go env GOPATH)/bin/git-repository

# Or on Windows
del %GOPATH%\bin\git-repository.exe
```

### Complete Cleanup

```bash
# Remove binary
sudo rm /usr/local/bin/git-repository

# Remove from GOPATH if installed there
rm $(go env GOPATH)/bin/git-repository

# Clear Go module cache (optional)
go clean -modcache
```

## Requirements

- **Go 1.21+** (only for building from source)
- **SSH keys** configured for private repositories
- **Git** is NOT required (uses go-git library)

## Examples

### Backup Your Dotfiles

```bash
git-repository save https://github.com/username/dotfiles ~/backups
# Creates: ~/backups/dotfiles.zip
```

### Backup Private Work Repository

```bash
git-repository save git@github.com:company/private-repo.git .
# Creates: ./private-repo.zip
```

### Automated Backups Script

```bash
#!/bin/bash
REPOS=(
  "https://github.com/user/repo1"
  "git@github.com:user/repo2.git"
)

for repo in "${REPOS[@]}"; do
  git-repository save "$repo" ~/backups/
done
```

## Troubleshooting

### "Permission denied" with SSH

Make sure your SSH key is added to your GitHub/GitLab account:
```bash
ssh-add ~/.ssh/id_rsa
ssh -T git@github.com  # Test connection
```

### "Command not found"

Add to your PATH:
```bash
export PATH="$PATH:/usr/local/bin"
# Add to ~/.bashrc or ~/.zshrc to make permanent
```

### Version shows "unknown"

The binary wasn't built with version information. Download from the [releases page](https://github.com/eeeelya/repository-backup/releases/latest) or build from source with proper flags.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

