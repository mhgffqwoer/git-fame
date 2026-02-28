package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

// CommitsSortType sorts data based on the number of commits.
type CommitsSortType struct{}

// NewCommitsSortType creates a new CommitsSortType instance.
func NewCommitsSortType() *CommitsSortType {
	return &CommitsSortType{}
}

// Compare compares two Data instances based on the number of commits.
func (s *CommitsSortType) Compare(a, b *metrics.Data) bool {
	if a.CommitsCount() != b.CommitsCount() {
		return a.CommitsCount() > b.CommitsCount()
	}
	if a.LinesCount() != b.LinesCount() {
		return a.LinesCount() > b.LinesCount()
	}
	if a.FilesCount() != b.FilesCount() {
		return a.FilesCount() > b.FilesCount()
	}
	return a.Name() < b.Name()
}
