package git_client

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/mhgffqwoer/git-fame/internal/shared"
)

type PathSetter interface {
	WithPath(path string) CommitSetter
}

type CommitSetter interface {
	WithCommit(commit string) Builder
}

type Builder interface {
	Build() (*GitClient, error)
}

type GitClientBuilder struct {
	path   string
	commit string
}

func NewBuilder() PathSetter {
	return &GitClientBuilder{}
}

func (b *GitClientBuilder) WithPath(path string) CommitSetter {
	b.path = path
	return b
}

func (b *GitClientBuilder) WithCommit(commit string) Builder {
	b.commit = commit
	return b
}

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
