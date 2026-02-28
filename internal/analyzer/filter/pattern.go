package filter

import (
	"path/filepath"

	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
)

// pattern represents a string pattern used for matching (either an extension or a glob).
type pattern string

func newFilterPattern(p string) pattern {
	return pattern(p)
}

// match checks if the given file's extension or path matches the pattern.
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
