package git_adapters

import (
	"context"
	"strings"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
	"github.com/mhgffqwoer/git-fame/internal/shared"
)

type HistoryReaderClient interface {
	FileHistory(path string) ([]byte, error)
	FileLastInfo(path string) ([]byte, error)
}

type FileHistoryStreamer struct {
	client    HistoryReaderClient
	committer bool
}

func NewFileHistoryStreamer(client HistoryReaderClient, committer bool) *FileHistoryStreamer {
	return &FileHistoryStreamer{
		client:    client,
		committer: committer,
	}
}

func (b *FileHistoryStreamer) StreamHistory(ctx context.Context, file *entity.File) (<-chan *entity.History, <-chan error) {
	lineCh := make(chan *entity.History)
	errCh := make(chan error, 1)

	go func() {
		defer close(lineCh)
		defer close(errCh)

		data, err := b.client.FileHistory(file.Path)
		if err != nil {
			errCh <- err
			return
		}

		if len(data) == 0 {
			lastInfo, err := b.client.FileLastInfo(file.Path)
			if err != nil {
				errCh <- err
				return
			}

			lines := strings.Split(string(lastInfo), "\n")
			if len(lines) >= 2 {
				commit := strings.TrimPrefix(lines[0], "commit ")
				author := strings.TrimPrefix(lines[1], "author ")
				hist := entity.NewHistory(file, commit, author, 0)
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case lineCh <- hist:
				}
			}
			return
		}

		lines := strings.Split(string(data), "\n")
		lines = shared.RemoveLinesStartingWith(lines,
			"author-mail", "author-time", "author-tz",
			"committer-mail", "committer-time", "committer-tz",
			"summary", "previous", "filename", "boundary", "\t",
		)
		blocks := shared.SplitIntoBlocks(lines, 3)

		for _, block := range blocks {
			var history *entity.History
			if b.committer {
				history = entity.NewHistory(
					file,
					strings.Fields(block[0])[0],
					strings.TrimPrefix(block[2], "committer "),
					1,
				)
			} else {
				history = entity.NewHistory(
					file,
					strings.Fields(block[0])[0],
					strings.TrimPrefix(block[1], "author "),
					1,
				)
			}

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case lineCh <- history:
			}
		}
	}()

	return lineCh, errCh
}
