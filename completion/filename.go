package completion

import (
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

type Filename string

var _ flags.Completer = Filename("")

func (Filename) Complete(match string) (output []flags.Completion) {
	names := listFilesByPattern("", match)
	result := make([]flags.Completion, 0, len(names))
	for _, name := range names {
		result = append(result, flags.Completion{Item: name})
	}
	return result
}

// ======================================================

func computeSearchWildcard(dir, match string) string {
	if len(dir) == 0 {
		return match + "*"
	}
	return dir + "/" + match + "*"
}

func listFilesByPattern(dir string, match string) []string {
	dir = filepath.Clean(dir)

	nameList, _ := filepath.Glob(computeSearchWildcard(dir, match))
	result := make([]string, 0, len(nameList))

	for _, name := range nameList {
		statInfo, _ := os.Stat(name)

		relativePath, err := filepath.Rel(dir, name)
		if err != nil {
			continue
		}

		if statInfo != nil {
			if statInfo.IsDir() {
				relativePath = relativePath + "/" + NoSpace
			}
		}
		result = append(result, relativePath)
	}

	return result
}

var globalListFilesByPatternFunc = listFilesByPattern
