package entity

// History represents the authorship history of a file.
type History struct {
	File       *File  // The file associated with this history entry.
	CommitHash string // The hash of the commit that last modified the file.
	Name       string // The name of the author who last modified the file.
	Lines      int    // The number of lines attributed to the author in this commit.
}

// String returns a string representation of the History entry.
func NewHistory(file *File, commitHash, name string, linesCnt int) *History {
	return &History{
		File:       file,
		CommitHash: commitHash,
		Name:       name,
		Lines:      linesCnt,
	}
}
