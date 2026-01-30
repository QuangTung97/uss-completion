package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeSearchWildcard(t *testing.T) {
	t.Run("empty dir", func(t *testing.T) {
		assert.Equal(t, "*", computeSearchWildcard("", ""))
		assert.Equal(t, "file01*", computeSearchWildcard("", "file01"))
	})

	t.Run("with dir", func(t *testing.T) {
		assert.Equal(t, "/path/dir01/*", computeSearchWildcard("/path/dir01", ""))
		assert.Equal(t, "/path/dir01/file01*", computeSearchWildcard("/path/dir01", "file01"))
	})
}
