package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

type FilesSortType struct{}

func NewFilesSortType() *FilesSortType {
	return &FilesSortType{}
}

func (s *FilesSortType) Compare(a, b *metrics.Data) bool {
	if a.FilesCount() != b.FilesCount() {
		return a.FilesCount() > b.FilesCount()
	}
	if a.LinesCount() != b.LinesCount() {
		return a.LinesCount() > b.LinesCount()
	}
	if a.CommitsCount() != b.CommitsCount() {
		return a.CommitsCount() > b.CommitsCount()
	}
	return a.Name() < b.Name()
}
