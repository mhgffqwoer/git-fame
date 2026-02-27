package filter

import (
	"path/filepath"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
)

type pattern string

func newFilterPattern(p string) pattern {
	return pattern(p)
}

func (f *pattern) match(file *entity.File) bool {
	var matched = false
	if file.Extension() == string(*f) {
		matched = true
	}
	temp, err := filepath.Match(string(*f), file.Path)
	if err != nil {
		return true
	}
	return matched || temp
}
