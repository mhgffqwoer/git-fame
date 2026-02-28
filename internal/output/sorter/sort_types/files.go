package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

// FilesSortType sorts data based on the number of files.
type FilesSortType struct{}

// NewFilesSortType creates a new FilesSortType instance.
func NewFilesSortType() *FilesSortType {
	return &FilesSortType{}
}

// Compare compares two Data instances based on the number of files.
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
