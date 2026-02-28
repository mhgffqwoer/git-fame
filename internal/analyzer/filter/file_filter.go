package filter

import (
	entity "github.com/mhgffqwoer/git-fame/internal/analyzer/entities"
)

// FileFilter determines if a specific file should be analyzed based on pre-defined patterns.
type FileFilter struct {
	includePatterns []pattern // Patterns that files must match to be included.
	excludePatterns []pattern // Patterns that files must NOT match to be included.
}

// newFileFilter is a constructor used by the builder to create a filter.
func newFileFilter(includePatterns, excludePatterns []pattern) *FileFilter {
	return &FileFilter{
		includePatterns: includePatterns,
		excludePatterns: excludePatterns,
	}
}

// include checks if the file matches at least one inclusion pattern (or if no inclusion patterns are set).
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

// exclude returns true if the file does NOT match any exclusion patterns.
func (f *FileFilter) exclude(file *entity.File) bool {
	for _, pattern := range f.excludePatterns {
		if pattern.match(file) {
			return false
		}
	}
	return true
}

// Allow returns true if the file matches at least one inclusion pattern and does NOT match any exclusion patterns.
func (f *FileFilter) Allow(file *entity.File) bool {
	return f.include(file) && f.exclude(file)
}
