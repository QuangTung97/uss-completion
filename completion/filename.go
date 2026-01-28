package completion

import "github.com/jessevdk/go-flags"

type Filename string

var _ flags.Completer = Filename("")

func (Filename) Complete(match string) (output []flags.Completion) {
	names := listFilesByPattern(match)
	result := make([]flags.Completion, 0, len(names))
	for _, name := range names {
		result = append(result, flags.Completion{Item: name})
	}
	return result
}
