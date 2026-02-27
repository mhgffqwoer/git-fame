package storage

import (
	"sync"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
	"github.com/mhgffqwoer/git-fame/internal/metrics"
)

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{
		m: sync.Map{},
	}
}

func (s *Storage) Save(history *entity.History) {
	userinfo := metrics.NewData(history.Name)
	if load, ok := s.m.Load(history.Name); ok {
		userinfo = load.(*metrics.Data)
	}
	userinfo.AddCommit(history.CommitHash)
	userinfo.AddFile(history.File.Path)
	userinfo.AddLines(history.Lines)
	s.m.Store(history.Name, userinfo)
}

func (s *Storage) Result() []*metrics.Data {
	result := make([]*metrics.Data, 0)
	s.m.Range(func(key, value any) bool {
		result = append(result, value.(*metrics.Data))
		return true
	})
	return result
}
