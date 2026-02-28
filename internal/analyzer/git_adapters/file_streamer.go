package git_adapters

import (
	"context"
	"strings"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
)

// FileListerClient defines the interface for listing files in a git repository.
type FileListerClient interface {
	ListFiles() ([]byte, error)
}

// FileStreamer defines the interface for retrieving a stream of all files in a repository.
type FileStreamer struct {
	client FileListerClient
}

// NewFileStreamer creates a new FileStreamer instance.
func NewFileStreamer(client FileListerClient) *FileStreamer {
	return &FileStreamer{client: client}
}

// StreamFiles retrieves a stream of all files in the repository.
func (f *FileStreamer) StreamFiles(ctx context.Context) (files <-chan *entity.File, errs <-chan error) {
	fileCh := make(chan *entity.File)
	errCh := make(chan error, 1)

	go func() {
		defer close(fileCh)
		defer close(errCh)

		out, err := f.client.ListFiles()
		if err != nil {
			errCh <- err
			return
		}

		list := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, path := range list {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fileCh <- entity.NewFile(path):
			}
		}
	}()

	return fileCh, errCh
}
