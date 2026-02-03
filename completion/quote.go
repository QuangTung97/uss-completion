package completion

import "strings"

const (
	ussPrefixConst  = "uss://"
	ussPrefixWithAt = "@" + ussPrefixConst
)

type QuoteHandler struct {
	withOpenQuote bool
	quoteCh       string

	prefixWithAt bool
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

func (h *QuoteHandler) equalUssPrefix(match string) bool {
	if match == ussPrefixConst {
		return true
	}
	if match == ussPrefixWithAt {
		h.prefixWithAt = true
		return true
	}
	return false
}

func (h *QuoteHandler) isPrefixOfUss(match string) bool {
	if strings.HasPrefix(ussPrefixConst, match) {
		return true
	}
	// TODO
	return false
}

func (h *QuoteHandler) getUssPrefix() string {
	if h.prefixWithAt {
		return ussPrefixWithAt
	}
	return ussPrefixConst
}

func (h *QuoteHandler) hasUssPrefix(match string) bool {
	if strings.HasPrefix(match, ussPrefixConst) {
		return true
	}
	if strings.HasPrefix(match, ussPrefixWithAt) {
		h.prefixWithAt = true
		return true
	}
	return false
}
