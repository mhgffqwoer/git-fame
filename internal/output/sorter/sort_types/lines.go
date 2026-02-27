package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

type LinesSortType struct{}

func NewLinesSortType() *LinesSortType {
	return &LinesSortType{}
}

func (s *LinesSortType) Compare(a, b *metrics.Data) bool {
	if a.LinesCount() != b.LinesCount() {
		return a.LinesCount() > b.LinesCount()
	}
	if a.CommitsCount() != b.CommitsCount() {
		return a.CommitsCount() > b.CommitsCount()
	}
	if a.FilesCount() != b.FilesCount() {
		return a.FilesCount() > b.FilesCount()
	}
	return a.Name() < b.Name()
}
