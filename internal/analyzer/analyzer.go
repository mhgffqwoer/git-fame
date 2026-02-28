package analyzer

import (
	"context"
	"fmt"
	"sync"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
	metrics "github.com/mhgffqwoer/git-fame/internal/metrics"
)

type FileStreamer interface {
	StreamFiles(ctx context.Context) (<-chan *entity.File, <-chan error)
}

type FileHistoryStreamer interface {
	StreamHistory(ctx context.Context, file *entity.File) (<-chan *entity.History, <-chan error)
}

type FileFilter interface {
	Allow(file *entity.File) bool
}

type Storage interface {
	Save(history *entity.History)
	Result() []*metrics.Data
}

type Analyzer struct {
	fileStreamer FileStreamer
	blameReader  FileHistoryStreamer
	filter       FileFilter
	storage      Storage
}

func NewAnalyzer(fileStreamer FileStreamer, blameReader FileHistoryStreamer, filter FileFilter, storage Storage) *Analyzer {
	return &Analyzer{
		fileStreamer: fileStreamer,
		blameReader:  blameReader,
		filter:       filter,
		storage:      storage,
	}
}

func (a *Analyzer) runFileStreamer(ctx context.Context) (files <-chan *entity.File, errs <-chan error) {
	return a.fileStreamer.StreamFiles(ctx)
}

func (a *Analyzer) runBlameReader(ctx context.Context, fileCh <-chan *entity.File) (histories <-chan *entity.History, errs <-chan error) {
	lineCh := make(chan *entity.History)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileCh {
				a.processBlameFile(ctx, file, lineCh, errCh)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(lineCh)
		close(errCh)
	}()
	return lineCh, errCh
}

func (a *Analyzer) processBlameFile(ctx context.Context, file *entity.File, lineCh chan<- *entity.History, errCh chan<- error) {
	fLineCh, fErrCh := a.blameReader.StreamHistory(ctx, file)

	for {
		select {
		case line, ok := <-fLineCh:
			if !ok {
				fLineCh = nil
			} else {
				lineCh <- line
			}
		case err, ok := <-fErrCh:
			if !ok {
				fErrCh = nil
			} else if err != nil {
				errCh <- fmt.Errorf("ошибка при обработке файла %s: %w", file.Path, err)
			}
		}
		if fLineCh == nil && fErrCh == nil {
			return
		}
	}
}

func (a *Analyzer) saveInfos(ctx context.Context, blameLine <-chan *entity.History) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case line, ok := <-blameLine:
				if !ok {
					return
				}
				a.storage.Save(line)
			case <-ctx.Done():
				return
			}
		}
	}()
	return done
}

func (a *Analyzer) runFileFilter(ctx context.Context, fileCh <-chan *entity.File) (files <-chan *entity.File, errs <-chan error) {
	fileChOut := make(chan *entity.File)
	errCh := make(chan error, 1)
	go func() {
		defer close(fileChOut)
		defer close(errCh)

		for file := range fileCh {
			if !a.filter.Allow(file) {
				continue
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fileChOut <- file:
			}
		}
	}()
	return fileChOut, errCh
}

func mergeErrors(errorChannels ...<-chan error) <-chan error {
	merged := make(chan error, 1)
	if len(errorChannels) == 0 {
		close(merged)
		return merged
	}
	var wg sync.WaitGroup
	for _, errChan := range errorChannels {
		wg.Add(1)
		go func(errChan <-chan error) {
			defer wg.Done()
			if err, ok := <-errChan; ok {
				select {
				case merged <- err:
				default:
				}
			}
		}(errChan)
	}
	go func() {
		wg.Wait()
		close(merged)
	}()
	return merged
}

func (a *Analyzer) Analyze(ctx context.Context) ([]*metrics.Data, error) {
	fileCh, fileErrCh := a.runFileStreamer(ctx)
	filterCh, filterErrCh := a.runFileFilter(ctx, fileCh)
	lineCh, lineErrCh := a.runBlameReader(ctx, filterCh)
	errCh := mergeErrors(lineErrCh, fileErrCh, filterErrCh)
	doneCh := a.saveInfos(ctx, lineCh)
	for {
		select {
		case <-doneCh:
			return a.storage.Result(), nil
		case err, ok := <-errCh:
			if ok && err != nil {
				return nil, err
			}
		}
	}
}
