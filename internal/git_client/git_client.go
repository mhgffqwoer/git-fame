package git_client

import (
	"bytes"
	"os/exec"
)

// GitClient represents a Git client.
type GitClient struct {
	absPath string
	commit  string
}

// newGitClient creates a new GitClient.
func newGitClient(path, commit string) *GitClient {
	return &GitClient{
		absPath: path,
		commit:  commit,
	}
}

// ListFiles lists all files in the repository.
func (g *GitClient) ListFiles() ([]byte, error) {
	checkCmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	checkCmd.Dir = g.absPath
	if err := checkCmd.Run(); err != nil {
		if g.commit == "HEAD" {
			return []byte{}, nil
		}
	}

	cmd := exec.Command("git", "ls-tree", "-r", "--name-only")
	cmd.Args = append(cmd.Args, g.commit)
	cmd.Dir = g.absPath

	return cmd.Output()
}

// FileHistory returns the history of a file.
func (g *GitClient) FileHistory(path string) ([]byte, error) {
	cmd := exec.Command("git", "blame", "--line-porcelain")
	cmd.Args = append(cmd.Args, g.commit, "--", path)
	cmd.Dir = g.absPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// FileLastInfo returns the last commit info of a file.
func (g *GitClient) FileLastInfo(path string) ([]byte, error) {
	cmd := exec.Command("git", "log", "--pretty=format:commit %H%nauthor %an")
	cmd.Args = append(cmd.Args, g.commit, "--", path)
	cmd.Dir = g.absPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
