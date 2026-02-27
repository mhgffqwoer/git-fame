package sort_types

import "github.com/mhgffqwoer/git-fame/internal/metrics"

type CommitsSortType struct{}

func NewCommitsSortType() *CommitsSortType {
	return &CommitsSortType{}
}

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
