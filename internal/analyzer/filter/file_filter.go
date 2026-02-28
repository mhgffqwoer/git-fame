package filter

import (
	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
)

type FileFilter struct {
	includePatterns []pattern
	excludePatterns []pattern
}

func newFileFilter(includePatterns, excludePatterns []pattern) *FileFilter {
	return &FileFilter{
		includePatterns: includePatterns,
		excludePatterns: excludePatterns,
	}
}

func (f *FileFilter) include(file *entity.File) bool {
	if len(f.includePatterns) == 0 {
		return true
	}
	for _, pattern := range f.includePatterns {
		if pattern.match(file) {
			return true
		}
	}
	return false
}

func (f *FileFilter) exclude(file *entity.File) bool {
	for _, pattern := range f.excludePatterns {
		if pattern.match(file) {
			return false
		}
	}
	return true
}

func (f *FileFilter) Allow(file *entity.File) bool {
	return f.include(file) && f.exclude(file)
}
