package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

// LinesSortType sorts data based on the number of lines.
type LinesSortType struct{}

// NewLinesSortType creates a new LinesSortType instance.
func NewLinesSortType() *LinesSortType {
	return &LinesSortType{}
}

// Compare compares two Data instances based on the number of lines.
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
