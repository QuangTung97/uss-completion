package completion

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

const limitListDataset = 20

func isZshShell() bool {
	return os.Getenv("GO_FLAGS_SHELL") == "zsh"
}

var IsZshShellFunc = isZshShell

// ------------------------------

type Uri string

var _ flags.Completer = Uri("")

func (Uri) Complete(match string) (output []flags.Completion) {
	return handleComplete(match, false)
}

func (u Uri) String() string {
	return string(u)
}

// ------------------------------

type UriAndFile string

var _ flags.Completer = UriAndFile("")

func (UriAndFile) Complete(match string) (output []flags.Completion) {
	return handleComplete(match, true)
}

func (u UriAndFile) String() string {
	return string(u)
}

// ------------------------------

type FileOrUri string

var _ flags.Completer = FileOrUri("")

func (FileOrUri) Complete(match string) (output []flags.Completion) {
	var empty Filename
	fileItems := empty.Complete(match)

	uriItems := handleComplete(match, true)
	fileItems = append(fileItems, uriItems...)

	return fileItems
}

func (u FileOrUri) String() string {
	return string(u)
}

// ------------------------------

func handleComplete(match string, withFile bool) (output []flags.Completion) {
	WriteToLog("Match: '%s'\n", match)
	defer func() {
		WriteToLog("Output: '%+v'\n", output)
	}()

	quote := &QuoteHandler{}

	if len(match) == 0 {
		return []flags.Completion{
			{Item: quote.getQuoteChar() + quote.getUssPrefix() + NoSpace},
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
	if quote.equalUssPrefix(match) {
		return nil
	}

	// match is prefix
	if quote.isPrefixOfUss(match) {
		return []flags.Completion{
			{Item: quote.getQuoteChar() + quote.getUssPrefix() + NoSpace},
		}
	}

	if !quote.hasUssPrefix(match) {
		return nil
	}

	openIndex := strings.Index(match, "{")
	closeIndex := strings.Index(match, "}")

	ussPrefix := quote.getUssPrefix()
	if openIndex <= 0 {
		matchDatasetName := match[len(ussPrefix):]
		datasetNames := GetMatchDatasetNamesFunc(matchDatasetName)

		prefix := quote.getQuoteChar() + ussPrefix

		var result []flags.Completion
		for _, name := range datasetNames {
			if strings.HasPrefix(name, matchDatasetName) {
				result = append(result, flags.Completion{
					Item: prefix + name + "{" + NoSpace,
				})
			}
		}

		if len(datasetNames) >= limitListDataset {
			result = append(result, flags.Completion{
				Item: prefix + matchDatasetName + NoSpace,
			})
		}

		return result
	}

	beforeBracketPart := match[:openIndex+1]
	if IsZshShellFunc() && !quote.withOpenQuote {
		beforeBracketPart = match[:openIndex]
	}

	if closeIndex <= 0 {
		prefix := quote.getQuoteChar() + beforeBracketPart
		attrsStr := match[openIndex+1:]
		datasetName := match[len(ussPrefix):openIndex]
		return handleAttrComplete(quote, prefix, datasetName, attrsStr)
	}

	prefix := quote.getQuoteChar() + match[:closeIndex+1] + quote.getQuoteChar()
	if IsZshShellFunc() {
		prefix = match[:closeIndex+1]
	}

	if !withFile {
		return []flags.Completion{
			{Item: prefix + NoSpace},
		}
	}

	remainMatch := match[closeIndex+1:]
	uriString := match[:closeIndex+1]
	searchDir := GetUriDiskPathFunc(uriString)

	// add no file suffix
	if len(remainMatch) == 0 {
		var result []flags.Completion
		result = append(result, flags.Completion{
			Item: prefix,
		})
		for _, fileMatch := range searchFilesWithNullDir(searchDir, "") {
			result = append(result, flags.Completion{
				Item: prefix + "/" + fileMatch,
			})
		}
		if len(result) == 1 {
			result[0].Item += NoSpace
		}
		return result
	}

	if !strings.HasPrefix(remainMatch, "/") {
		return nil
	}

	// remove prefix
	remainMatch = remainMatch[1:]

	var result []flags.Completion
	for _, fileMatch := range searchFilesWithNullDir(searchDir, remainMatch) {
		if !strings.HasPrefix(fileMatch, remainMatch) {
			continue
		}
		result = append(result, flags.Completion{
			Item: prefix + "/" + fileMatch,
		})
	}
	return result
}

func searchFilesWithNullDir(searchDir string, remainMatch string) []string {
	if searchDir == NullDir {
		return nil
	}
	return globalListFilesByPatternFunc(searchDir, remainMatch)
}

const attrDateKey = "date"

func handleAttrComplete(
	quote *QuoteHandler, prefix string, datasetName string, attrsStr string,
) []flags.Completion {
	attrsStr = strings.TrimSpace(attrsStr)

	lastAttr := attrsStr
	lastCommaIndex := strings.LastIndex(attrsStr, ",")
	existedKeys := map[string]struct{}{}
	if lastCommaIndex >= 0 {
		lastAttr = strings.TrimSpace(attrsStr[lastCommaIndex+1:])
		if !IsZshShellFunc() || quote.withOpenQuote {
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
		attrDateKey: {
			"date=",
		},
	}

	keyList := []string{attrDateKey}
	keySet := map[string]struct{}{
		attrDateKey: {},
	}

	versionList := GetAllVersionsFunc(datasetName)
	if versionList.DisableCompletion {
		return nil
	}

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

	isLastKey := len(existedKeys) >= len(keyList)-1

	var result []flags.Completion
	for _, attrKey := range keyList {
		_, existed := existedKeys[attrKey]
		if existed {
			continue
		}

		handleDateAttrComplete(quote, attrKey, lastAttr, isLastKey, prefix, &result)

		for _, kv := range allMatches[attrKey] {
			if !strings.HasPrefix(kv, lastAttr) {
				continue
			}

			matchStr := kv
			if attrKey != attrDateKey && isLastKey {
				matchStr += "}" + quote.getQuoteChar()
			}

			result = append(result, flags.Completion{
				Item: prefix + matchStr + NoSpace,
			})
		}
	}
	return result
}

func handleDateAttrComplete(
	quote *QuoteHandler,
	attrKey string, lastAttr string,
	isLastKey bool, prefix string,
	result *[]flags.Completion,
) {
	if attrKey != attrDateKey {
		return
	}

	if !strings.HasPrefix(lastAttr, attrDateKey) {
		return
	}

	equalIndex := strings.Index(lastAttr, "=")
	if equalIndex <= 0 {
		return
	}

	value := strings.TrimSpace(lastAttr[equalIndex+1:])
	if len(value) < 8 {
		return
	}

	if isLastKey {
		*result = append(*result, flags.Completion{
			Item: prefix + lastAttr + "}" + quote.getQuoteChar() + NoSpace,
		})
	} else {
		*result = append(*result, flags.Completion{
			Item: prefix + lastAttr + "," + NoSpace,
		})
	}
}
