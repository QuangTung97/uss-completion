package completion

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Uri string

var _ flags.Completer = Uri("")

func (Uri) Complete(match string) (output []flags.Completion) {
	return handleComplete(match, false)
}

// ------------------------------

type UriAndFile string

var _ flags.Completer = UriAndFile("")

func (UriAndFile) Complete(match string) (output []flags.Completion) {
	return handleComplete(match, true)
}

// ======================================================

func listFilesByPattern(match string) []string {
	nameList, _ := filepath.Glob(match + "*")
	result := make([]string, 0, len(nameList))

	for _, name := range nameList {
		statInfo, _ := os.Stat(name)
		if statInfo != nil {
			if statInfo.IsDir() {
				name = name + "/" + NoSpace
			}
		}

		result = append(result, name)
	}

	return result
}

var globalListFilesByPatternFunc = listFilesByPattern

func removeQuoted(match string, withOpenQuote *bool) string {
	if !strings.HasPrefix(match, DoubleQuote) {
		return match
	}

	*withOpenQuote = true

	match = strings.TrimPrefix(match, DoubleQuote)
	closeIndex := strings.Index(match, DoubleQuote)
	if closeIndex < 0 {
		return match
	}
	return match[:closeIndex] + match[closeIndex+1:]
}

func handleComplete(match string, withFile bool) (output []flags.Completion) {
	WriteToLog("Match: '%s'\n", match)
	defer func() {
		WriteToLog("Output: '%+v'\n", output)
	}()

	var withOpenQuote bool
	match = removeQuoted(match, &withOpenQuote)

	const ussPrefix = "uss://"
	if match == ussPrefix {
		return nil
	}

	// match is prefix
	if strings.HasPrefix(ussPrefix, match) {
		prefix := DoubleQuote
		return []flags.Completion{
			{Item: prefix + `uss://` + NoSpace},
		}
	}

	if !strings.HasPrefix(match, ussPrefix) {
		return nil
	}

	openIndex := strings.Index(match, "{")
	closeIndex := strings.Index(match, "}")

	var resultWithOpenQuote []flags.Completion
	if !withOpenQuote {
		resultWithOpenQuote = append(resultWithOpenQuote, flags.Completion{Item: DoubleQuote + match})
	}

	if openIndex <= 0 {
		return resultWithOpenQuote
	}
	if closeIndex <= 0 {
		return resultWithOpenQuote
	}

	prefix := DoubleQuote + match[:closeIndex+1] + DoubleQuote
	if !withFile {
		return []flags.Completion{
			{Item: prefix},
		}
	}

	remainMatch := match[closeIndex+1:]
	var result []flags.Completion

	// add no file suffix
	if len(remainMatch) == 0 {
		result = append(result, flags.Completion{
			Item: prefix,
		})
		for _, fileMatch := range globalListFilesByPatternFunc(remainMatch) {
			result = append(result, flags.Completion{
				Item: prefix + "/" + fileMatch,
			})
		}
		return result
	}

	if !strings.HasPrefix(remainMatch, "/") {
		return nil
	}

	// remove prefix
	remainMatch = remainMatch[1:]

	for _, fileMatch := range globalListFilesByPatternFunc(remainMatch) {
		if !strings.HasPrefix(fileMatch, remainMatch) {
			continue
		}
		result = append(result, flags.Completion{
			Item: prefix + "/" + fileMatch,
		})
	}
	return result
}
