package completion

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

const ussPrefix = "uss://"

func isZshShell() bool {
	return os.Getenv("GO_FLAGS_SHELL") == "zsh"
}

var IsZshShellFunc = isZshShell

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

func handleComplete(match string, withFile bool) (output []flags.Completion) {
	WriteToLog("Match: '%s'\n", match)
	defer func() {
		WriteToLog("Output: '%+v'\n", output)
	}()

	quote := &QuoteHandler{}

	if len(match) == 0 {
		return []flags.Completion{
			{Item: quote.getQuoteChar() + ussPrefix + NoSpace},
		}
	}

	output = doHandleComplete(quote, match, withFile)
	if IsZshShellFunc() {
		for i := range output {
			item := output[i].Item
			output[i].Item = strings.TrimPrefix(item, quote.getQuoteChar())
		}
	}
	return output
}

func doHandleComplete(quote *QuoteHandler, match string, withFile bool) []flags.Completion {
	output := coreHandleComplete(quote, match, withFile)
	if len(output) != 1 {
		return output
	}

	if !quote.withOpenQuote {
		return output
	}

	item := output[0].Item
	if !strings.HasSuffix(item, NoSpace) {
		return output
	}

	item = strings.TrimSuffix(item, NoSpace)
	if item == match {
		return nil
	}

	if IsZshShellFunc() {
		return output
	}

	return []flags.Completion{
		{Item: item + BlackBullet + NoSpace},
		{Item: item + WhiteBullet + NoSpace},
	}
}

func coreHandleComplete(
	quote *QuoteHandler, match string, withFile bool,
) []flags.Completion {
	match = quote.removeQuoted(match)
	if match == ussPrefix {
		return nil
	}

	// match is prefix
	if strings.HasPrefix(ussPrefix, match) {
		return []flags.Completion{
			{Item: quote.getQuoteChar() + ussPrefix + NoSpace},
		}
	}

	if !strings.HasPrefix(match, ussPrefix) {
		return nil
	}

	openIndex := strings.Index(match, "{")
	closeIndex := strings.Index(match, "}")

	var resultWithOpenQuote []flags.Completion
	if !quote.withOpenQuote {
		resultWithOpenQuote = append(resultWithOpenQuote, flags.Completion{
			Item: quote.getQuoteChar() + match,
		})
	}

	if openIndex <= 0 {
		return resultWithOpenQuote
	}

	beforeBracketPart := match[:openIndex+1]
	if IsZshShellFunc() && !quote.withOpenQuote {
		beforeBracketPart = match[:openIndex]
	}

	if closeIndex <= 0 {
		prefix := quote.getQuoteChar() + beforeBracketPart
		attrsStr := match[openIndex+1:]
		return handleAttrComplete(quote, prefix, attrsStr, quote.withOpenQuote)
	}

	prefix := quote.getQuoteChar() + match[:closeIndex+1] + quote.getQuoteChar()
	if IsZshShellFunc() {
		prefix = match[:closeIndex+1]
	}

	if !withFile {
		return []flags.Completion{
			{Item: prefix},
		}
	}

	remainMatch := match[closeIndex+1:]
	uriString := match[:closeIndex+1]
	searchDir := GetUriDiskPathFunc(uriString)

	var result []flags.Completion

	// add no file suffix
	if len(remainMatch) == 0 {
		result = append(result, flags.Completion{
			Item: prefix,
		})
		for _, fileMatch := range globalListFilesByPatternFunc(searchDir, "") {
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

	for _, fileMatch := range globalListFilesByPatternFunc(searchDir, remainMatch) {
		if !strings.HasPrefix(fileMatch, remainMatch) {
			continue
		}
		result = append(result, flags.Completion{
			Item: prefix + "/" + fileMatch,
		})
	}
	return result
}

func handleAttrComplete(
	quote *QuoteHandler,
	prefix string, attrsStr string, withOpenQuote bool,
) []flags.Completion {
	attrsStr = strings.TrimSpace(attrsStr)

	lastAttr := attrsStr
	lastCommaIndex := strings.LastIndex(attrsStr, ",")
	existedKeys := map[string]struct{}{}
	if lastCommaIndex >= 0 {
		lastAttr = strings.TrimSpace(attrsStr[lastCommaIndex+1:])
		if !IsZshShellFunc() || withOpenQuote {
			prefix = prefix + attrsStr[:lastCommaIndex+1]
		}

		for _, kv := range strings.Split(attrsStr[:lastCommaIndex], ",") {
			kv = strings.TrimSpace(kv)
			equalIndex := strings.Index(kv, "=")
			if equalIndex <= 0 {
				continue
			}
			key := kv[:equalIndex]
			existedKeys[key] = struct{}{}
		}
	}

	allMatches := map[string][]string{
		"date": {
			"date=",
		},
	}

	keyList := []string{"date"}
	keySet := map[string]struct{}{
		"date": {},
	}

	versionList := GetAllVersionsFunc()
	for _, version := range versionList.Versions {
		for key, val := range version {
			allMatches[key] = append(allMatches[key], key+"="+val)
			_, existed := keySet[key]
			if !existed {
				keyList = append(keyList, key)
				keySet[key] = struct{}{}
			}
		}
	}

	var result []flags.Completion
	for _, attrKey := range keyList {
		_, existed := existedKeys[attrKey]
		if existed {
			continue
		}

		for _, kv := range allMatches[attrKey] {
			if !strings.HasPrefix(kv, lastAttr) {
				continue
			}

			matchStr := kv
			if attrKey != "date" && len(existedKeys) >= len(keyList)-1 {
				matchStr += "}" + quote.getQuoteChar()
			}

			result = append(result, flags.Completion{
				Item: prefix + matchStr + NoSpace,
			})
		}
	}
	return result
}
