package git_client

import (
	"bytes"
	"os/exec"
)

type GitClient struct {
	absPath string
	commit  string
}

func newGitClient(path, commit string) *GitClient {
	return &GitClient{
		absPath: path,
		commit:  commit,
	}
}

func (g *GitClient) ListFiles() ([]byte, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only")
	cmd.Args = append(cmd.Args, g.commit)
	cmd.Dir = g.absPath
	return cmd.Output()
}

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
