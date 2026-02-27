package sorter

import (
	"fmt"
	"sort"

	"github.com/mhgffqwoer/git-fame/internal/metrics"
	"github.com/mhgffqwoer/git-fame/internal/output/sorter/sort_types"
)

type SortType interface {
	Compare(a, b *metrics.Data) bool
}

type Sorter struct {
	sortType SortType
}

func NewSorter(formatType string) (*Sorter, error) {
	switch formatType {
	case "lines":
		return &Sorter{
			sortType: sort_types.NewLinesSortType(),
		}, nil
	case "commits":
		return &Sorter{
			sortType: sort_types.NewCommitsSortType(),
		}, nil
	case "files":
		return &Sorter{
			sortType: sort_types.NewFilesSortType(),
		}, nil
	default:
		return nil, fmt.Errorf("error")
	}
}

func (f *Sorter) Sort(list []*metrics.Data) {
	sort.SliceStable(list, func(i, j int) bool {
		return f.sortType.Compare(list[i], list[j])
	})
}
