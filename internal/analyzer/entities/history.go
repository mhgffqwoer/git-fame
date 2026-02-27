package entity

type History struct {
	File       *File
	CommitHash string
	Name       string
	Lines      int
}

func NewHistory(file *File, commitHash, name string, linesCnt int) *History {
	return &History{
		File:       file,
		CommitHash: commitHash,
		Name:       name,
		Lines:      linesCnt,
	}
}
