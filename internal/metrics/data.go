package metrics

import "sync"

// Data holds the metrics for a single user.
type Data struct {
	mu      sync.Mutex          // Mutex to protect concurrent access to the Data fields.
	name    string              // The name of the user (author/committer).
	lines   int                 // Total number of lines attributed to the user.
	commits map[string]struct{} // Set of unique commit hashes attributed to the user.
	files   map[string]struct{} // Set of unique file paths attributed to the user.
}

// NewData creates a new Data object with the given name.
func NewData(name string) *Data {
	return &Data{
		name:    name,
		commits: make(map[string]struct{}),
		files:   make(map[string]struct{}),
	}
}

// AddCommit adds a new commit hash to the user's data.
func (u *Data) AddCommit(commitHash string) {
	u.mu.Lock()
	u.commits[commitHash] = struct{}{}
	u.mu.Unlock()
}

// AddFile adds a new file path to the user's data.
func (u *Data) AddFile(path string) {
	u.mu.Lock()
	u.files[path] = struct{}{}
	u.mu.Unlock()
}

// AddLines adds a number of lines to the user's total line count.
func (u *Data) AddLines(lines int) {
	u.mu.Lock()
	u.lines += lines
	u.mu.Unlock()
}

// Name returns the name of the user.
func (u *Data) Name() string {
	return u.name
}

// CommitsCount returns the number of unique commit hashes attributed to the user.
func (u *Data) CommitsCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.commits)
}

// FilesCount returns the number of unique file paths attributed to the user.
func (u *Data) FilesCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.files)
}

// LinesCount returns the total number of lines attributed to the user.
func (u *Data) LinesCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.lines
}
