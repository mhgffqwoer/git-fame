package git_adapters

import (
	"context"
	"strings"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
	"github.com/mhgffqwoer/git-fame/internal/shared"
)

// HistoryReaderClient defines required git operations for blame and metadata retrieval.
type HistoryReaderClient interface {
	FileHistory(path string) ([]byte, error)
	FileLastInfo(path string) ([]byte, error)
}

// FileHistoryStreamer handles the retrieval and parsing of file authorship history.
type FileHistoryStreamer struct {
	client    HistoryReaderClient
	committer bool
}

// NewFileHistoryStreamer creates a new streamer for file history.
func NewFileHistoryStreamer(client HistoryReaderClient, committer bool) *FileHistoryStreamer {
	return &FileHistoryStreamer{
		client:    client,
		committer: committer,
	}
}

// StreamHistory retrieves blame data for a file and streams individual history records.
func (b *FileHistoryStreamer) StreamHistory(ctx context.Context, file *entity.File) (histories <-chan *entity.History, errs <-chan error) {
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
			b.handleEmptyHistory(ctx, file, lineCh, errCh)
			return
		}

		b.processHistoryData(ctx, file, data, lineCh, errCh)
	}()

	return lineCh, errCh
}

// handleEmptyHistory handles the case where a file has no history.
func (b *FileHistoryStreamer) handleEmptyHistory(ctx context.Context, file *entity.File, lineCh chan<- *entity.History, errCh chan<- error) {
	lastInfo, err := b.client.FileLastInfo(file.Path)
	if err != nil {
		errCh <- err
		return
	}

	lines := strings.Split(string(lastInfo), "\n")
	if len(lines) < 2 {
		return
	}

	commit := strings.TrimPrefix(lines[0], "commit ")
	author := strings.TrimPrefix(lines[1], "author ")
	hist := entity.NewHistory(file, commit, author, 0)

	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case lineCh <- hist:
	}
}

// processHistoryData processes blame data and streams individual history records.
func (b *FileHistoryStreamer) processHistoryData(ctx context.Context, file *entity.File, data []byte, lineCh chan<- *entity.History, errCh chan<- error) {
	lines := strings.Split(string(data), "\n")
	lines = shared.RemoveLinesStartingWith(lines,
		"author-mail", "author-time", "author-tz",
		"committer-mail", "committer-time", "committer-tz",
		"summary", "previous", "filename", "boundary", "\t",
	)
	blocks := shared.SplitIntoBlocks(lines, 3)

	for _, block := range blocks {
		history := b.parseHistoryBlock(file, block)
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case lineCh <- history:
		}
	}
}

// parseHistoryBlock parses a blame block and returns a history record.
func (b *FileHistoryStreamer) parseHistoryBlock(file *entity.File, block []string) *entity.History {
	commit := strings.Fields(block[0])[0]
	if b.committer {
		return entity.NewHistory(file, commit, strings.TrimPrefix(block[2], "committer "), 1)
	}
	return entity.NewHistory(file, commit, strings.TrimPrefix(block[1], "author "), 1)
}
