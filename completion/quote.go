package completion

import "strings"

type QuoteHandler struct {
	withOpenQuote bool
	quoteCh       string
}

func (h *QuoteHandler) getQuoteChar() string {
	if len(h.quoteCh) > 0 {
		return h.quoteCh
	}
	return `"`
}

func (h *QuoteHandler) removeQuoted(match string) string {
	if len(match) == 0 {
		return match
	}

	if match[0] != '"' && match[0] != '\'' {
		return match
	}

	h.quoteCh = string(match[0])
	h.withOpenQuote = true
	match = match[1:]

	closeIndex := strings.Index(match, h.quoteCh)
	if closeIndex < 0 {
		return match
	}
	return match[:closeIndex] + match[closeIndex+1:]
}
