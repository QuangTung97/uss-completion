package completion

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type UriValue string

var _ flags.Completer = UriValue("")

func removeQuoted(match string) (string, bool, string) {
	if !strings.HasPrefix(match, DoubleQuote) {
		return match, false, match
	}

	match = strings.TrimPrefix(match, DoubleQuote)

	closeIndex := strings.Index(match, DoubleQuote)
	if closeIndex >= 0 {
		return match[:closeIndex], true, match[closeIndex+1:]
	}
	return match, false, match
}

func (v UriValue) Complete(match string) (output []flags.Completion) {
	WriteToLog("Match: '%s'\n", match)
	defer func() {
		WriteToLog("Output: '%+v'\n", output)
	}()

	match, fullyQuoted, remainMatch := removeQuoted(match)

	const ussPrefix = "uss://"
	if match == ussPrefix {
		return nil
	}

	// match is prefix
	if strings.HasPrefix(ussPrefix, match) {
		return []flags.Completion{
			{Item: DoubleQuote + `uss://` + NoSpace},
		}
	}

	if !strings.HasPrefix(match, ussPrefix) {
		return nil
	}

	openIndex := strings.Index(match, "{")
	closeIndex := strings.Index(match, "}")

	if openIndex > 0 && closeIndex > 0 {
		if fullyQuoted {
			prefix := DoubleQuote + match + DoubleQuote

			// TODO handle correctly
			var result []flags.Completion
			result = append(result, flags.Completion{
				Item: prefix,
			})

			empty := flags.Filename("")
			for _, fileMatch := range empty.Complete(remainMatch) {
				result = append(result, flags.Completion{
					Item: prefix + "/" + fileMatch.Item,
				})
			}
			return result
		}
		return []flags.Completion{
			{Item: DoubleQuote + match + DoubleQuote + NoSpace},
		}
	}

	return nil
}

func WriteToLog(format string, args ...any) {
	if os.Getenv("GO_TEST") != "" {
		return
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(file, format, args...)
	if err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}
