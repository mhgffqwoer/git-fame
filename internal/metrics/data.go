package metrics

import "sync"

type Data struct {
	mu      sync.Mutex
	name    string
	lines   int
	commits map[string]struct{}
	files   map[string]struct{}
}

func NewData(name string) *Data {
	return &Data{
		name:    name,
		commits: make(map[string]struct{}),
		files:   make(map[string]struct{}),
	}
}

func (u *Data) AddCommit(commitHash string) {
	u.mu.Lock()
	u.commits[commitHash] = struct{}{}
	u.mu.Unlock()
}

func (u *Data) AddFile(path string) {
	u.mu.Lock()
	u.files[path] = struct{}{}
	u.mu.Unlock()
}

func (u *Data) AddLines(lines int) {
	u.mu.Lock()
	u.lines += lines
	u.mu.Unlock()
}

func (u *Data) Name() string {
	return u.name
}

func (u *Data) CommitsCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.commits)
}

func (u *Data) FilesCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.files)
}

func (u *Data) LinesCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.lines
}
