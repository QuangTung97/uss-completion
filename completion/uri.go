package completion

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type UriValue string

var _ flags.Completer = UriValue("")

func listFilesByPattern(match string) []string {
	empty := flags.Filename("")

	items := empty.Complete(match)
	result := make([]string, 0, len(items))

	for _, it := range items {
		result = append(result, it.Item)
	}

	return result
}

var globalListFilesByPatternFunc = listFilesByPattern

func removeQuoted(match string, withClosedQuote *bool) string {
	if !strings.HasPrefix(match, DoubleQuote) {
		return match
	}

	match = strings.TrimPrefix(match, DoubleQuote)
	closeIndex := strings.Index(match, DoubleQuote)
	if closeIndex >= 0 {
		*withClosedQuote = true
		return match[:closeIndex] + match[closeIndex+1:]
	}
	return match
}

func (v UriValue) Complete(match string) (output []flags.Completion) {
	WriteToLog("Match: '%s'\n", match)
	defer func() {
		WriteToLog("Output: '%+v'\n", output)
	}()

	var withClosedQuote bool
	match = removeQuoted(match, &withClosedQuote)

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

	if openIndex <= 0 {
		return nil
	}
	if closeIndex <= 0 {
		return nil
	}

	remainMatch := match[closeIndex+1:]
	prefix := DoubleQuote + match[:closeIndex+1] + DoubleQuote
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

func PrintArray(values []string) string {
	var buf bytes.Buffer
	for index, v := range values {
		if index > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%q", v))
	}
	return buf.String()
}
