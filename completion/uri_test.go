package completion

import (
	"os"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func TestUriValue_Complete(t *testing.T) {
	_ = os.Setenv("GO_TEST", "1")

	t.Run("empty", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t, []flags.Completion{
			{Item: `"uss://<NS>`},
		}, v.Complete(""))
	})

	t.Run("with uss prefix", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t, []flags.Completion{
			{Item: DoubleQuote + "uss://<NS>"},
		}, v.Complete("us"))
	})

	t.Run("with full prefix", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t, []flags.Completion(nil), v.Complete(DoubleQuote+"uss://"))
	})

	t.Run("full", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t,
			[]flags.Completion{
				{Item: `"uss://test01{date=20250219}"<NS>`},
			},
			v.Complete("uss://test01{date=20250219}"),
		)
	})

	t.Run("near full with close bracket", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t,
			[]flags.Completion{
				{Item: `"uss://test01{date=20250219,asset_type=at}"<NS>`},
			},
			v.Complete(DoubleQuote+"uss://test01{date=20250219,asset_type=at}"),
		)
	})

	t.Run("with full quote", func(t *testing.T) {
		v := UriValue("")
		assert.Equal(t,
			[]flags.Completion(nil),
			v.Complete(`"uss://test01{date=20250219,asset_type=at}"`),
		)
	})
}
