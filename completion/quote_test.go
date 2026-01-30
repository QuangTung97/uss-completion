package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteHandler(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := &QuoteHandler{}
		s := h.removeQuoted(``)
		assert.Equal(t, ``, s)
		assert.Equal(t, false, h.withOpenQuote)
		assert.Equal(t, ``, h.quoteCh)
	})

	t.Run("normal", func(t *testing.T) {
		h := &QuoteHandler{}
		s := h.removeQuoted(`"hello world"`)
		assert.Equal(t, `hello world`, s)
		assert.Equal(t, true, h.withOpenQuote)
		assert.Equal(t, `"`, h.quoteCh)
	})

	t.Run("single quote", func(t *testing.T) {
		h := &QuoteHandler{}
		s := h.removeQuoted(`'hello world'`)
		assert.Equal(t, `hello world`, s)
		assert.Equal(t, true, h.withOpenQuote)
		assert.Equal(t, `'`, h.quoteCh)
	})

	t.Run("without quote", func(t *testing.T) {
		h := &QuoteHandler{}
		s := h.removeQuoted(`hello world`)
		assert.Equal(t, `hello world`, s)
		assert.Equal(t, false, h.withOpenQuote)
		assert.Equal(t, "", h.quoteCh)
	})
}
