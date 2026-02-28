package git_client

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/mhgffqwoer/git-fame/internal/shared"
)

// PathSetter defines the interface for setting the repository path in the builder.
type PathSetter interface {
	WithPath(path string) CommitSetter
}

// CommitSetter defines the interface for setting the commit hash in the builder.
type CommitSetter interface {
	WithCommit(commit string) Builder
}

// Builder defines the interface for building a GitClient.
type Builder interface {
	Build() (*GitClient, error)
}

// GitClientBuilder builds a GitClient.
type GitClientBuilder struct {
	path   string // Repository path.
	commit string // Commit hash to analyze (e.g., "HEAD" or a specific commit).
}

// NewBuilder creates a new GitClientBuilder.
func NewBuilder() PathSetter {
	return &GitClientBuilder{}
}

// WithPath sets the repository path in the builder.
func (b *GitClientBuilder) WithPath(path string) CommitSetter {
	b.path = path
	return b
}

// WithCommit sets the commit hash in the builder.
func (b *GitClientBuilder) WithCommit(commit string) Builder {
	b.commit = commit
	return b
}

// Build builds a GitClient.
func (b *GitClientBuilder) Build() (*GitClient, error) {
	abs, err := shared.ToAbsolutePath(b.path)
	if err != nil {
		return nil, err
	}
	if err := ensureGitRepository(abs); err != nil {
		return nil, err
	}
	return newGitClient(abs, b.commit), nil
}

// ensureGitRepository checks if a directory is a git repository.
func ensureGitRepository(absPath string) error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = absPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("directory %q is not a git repository: %w", absPath, err)
	}
	if out.String() != "true\n" {
		return fmt.Errorf("directory %q is not a git repository", absPath)
	}
	return nil
}
