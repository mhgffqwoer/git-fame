package shared

import (
	"os"
	"path/filepath"
	"strings"
)

func ToAbsolutePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func SplitIntoBlocks(lines []string, blockSize int) [][]string {
	if len(lines) == 0 {
		return nil
	}
	count := (len(lines) + blockSize - 1) / blockSize
	blocks := make([][]string, 0, count)
	for i := 0; i < len(lines); i += blockSize {
		end := min(i+blockSize, len(lines))
		block := lines[i:end]

		blocks = append(blocks, block)
	}
	if len(blocks) == 0 {
		return blocks
	}
	return blocks[:len(blocks)-1]
}

func RemoveLinesStartingWith(lines []string, words ...string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		flag := true
		for _, word := range words {
			if strings.HasPrefix(line, word) {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, line)
		}
	}
	return result
}
