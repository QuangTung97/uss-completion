package completion

import (
	"os"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func TestUriAndFile_Complete__Basic(t *testing.T) {
	_ = os.Setenv("GO_TEST", "1")

	t.Run("empty", func(t *testing.T) {
		v := UriAndFile("")
		assert.Equal(t, []flags.Completion{
			{Item: `"uss://<NS>`},
		}, v.Complete(""))
	})

	t.Run("with uss prefix", func(t *testing.T) {
		v := UriAndFile("")
		assert.Equal(t, []flags.Completion{
			{Item: `"uss://<NS>`},
		}, v.Complete("us"))
	})

	t.Run("with uss prefix, with open double quote", func(t *testing.T) {
		v := UriAndFile("")
		assert.Equal(t, []flags.Completion{
			{Item: `"uss://` + BlackBullet + NoSpace},
			{Item: `"uss://` + WhiteBullet + NoSpace},
		}, v.Complete(`"uss`))
	})

	t.Run("with full prefix", func(t *testing.T) {
		v := UriAndFile("")
		assert.Equal(t, []flags.Completion(nil), v.Complete(`"uss://`))
	})
}

type uriValueTest struct {
	fileMatchDir   string
	fileMatchInput string
	fileMatchCalls int
	fileList       []string
}

func newUriValueTest(t *testing.T) *uriValueTest {
	_ = os.Setenv("GO_TEST", "1")

	v := &uriValueTest{}

	prevFunc := globalListFilesByPatternFunc
	globalListFilesByPatternFunc = func(dir string, match string) []string {
		v.fileMatchDir = dir
		v.fileMatchInput = match
		v.fileMatchCalls++
		return v.fileList
	}
	t.Cleanup(func() {
		globalListFilesByPatternFunc = prevFunc
	})

	return v
}

func (v *uriValueTest) completeUriAndFile(prefix string) []string {
	value := UriAndFile("")
	items := value.Complete(prefix)
	result := make([]string, 0, len(items))
	for _, it := range items {
		result = append(result, it.Item)
	}
	return result
}

func (v *uriValueTest) completeUri(prefix string) []string {
	value := Uri("")
	items := value.Complete(prefix)
	result := make([]string, 0, len(items))
	for _, it := range items {
		result = append(result, it.Item)
	}
	return result
}

func TestUriAndFile_Complete__With_Files(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250219,asset_type=at}"`,
				`"uss://test01{date=20250219,asset_type=at}"/file01`,
				`"uss://test01{date=20250219,asset_type=at}"/file02`,
			},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at}"`),
		)

		// check input
		assert.Equal(t, "", v.fileMatchInput)
		assert.Equal(t, 1, v.fileMatchCalls)
	})

	t.Run("with dataset name, and no quote", func(t *testing.T) {
		v := newUriValueTest(t)

		assert.Equal(
			t,
			[]string{
				`"uss://test01`,
			},
			v.completeUriAndFile(`uss://test01`),
		)

		// check input
		assert.Equal(t, 0, v.fileMatchCalls)
	})

	t.Run("with match filename", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
			"example",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250219,asset_type=at}"/file01`,
				`"uss://test01{date=20250219,asset_type=at}"/file02`,
			},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at}"/file`),
		)

		// check input
		assert.Equal(t, "file", v.fileMatchInput)
	})

	t.Run("no close quote, with match filename", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
			"example",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250219,asset_type=at}"/file01`,
				`"uss://test01{date=20250219,asset_type=at}"/file02`,
			},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at}/file`),
		)

		// check input
		assert.Equal(t, "file", v.fileMatchInput)
	})

	t.Run("no close quote, with match filename, single quote", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
			"example",
		}

		assert.Equal(
			t,
			[]string{
				`'uss://test01{date=20250219,asset_type=at}'/file01`,
				`'uss://test01{date=20250219,asset_type=at}'/file02`,
			},
			v.completeUriAndFile(`'uss://test01{date=20250219,asset_type=at}/file`),
		)

		// check input
		assert.Equal(t, "file", v.fileMatchInput)
	})

	t.Run("no close quote, no match filename", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
			"example",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250219,asset_type=at}"`,
				`"uss://test01{date=20250219,asset_type=at}"/file01`,
				`"uss://test01{date=20250219,asset_type=at}"/file02`,
				`"uss://test01{date=20250219,asset_type=at}"/example`,
			},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at}`),
		)

		// check input
		assert.Equal(t, "", v.fileMatchInput)
	})

	t.Run("without close bracket", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at`),
		)
	})

	t.Run("without open bracket", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"uss://test01date=20250219,asset_type=at}`),
		)
	})

	t.Run("without slash", func(t *testing.T) {
		v := newUriValueTest(t)

		v.fileList = []string{
			"file01",
			"file02",
		}

		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"uss://test01{date=20250219,asset_type=at}file`),
		)

		// check input
		assert.Equal(t, "", v.fileMatchInput)
	})

	t.Run("with attr completion at start", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=<NS>`,
				`"uss://test01{asset_type=equity<NS>`,
				`"uss://test01{asset_type=options<NS>`,
			},
			v.completeUriAndFile(`"uss://test01{`),
		)
		assert.Equal(t, 0, v.fileMatchCalls)
	})

	t.Run("attr completion with prefix", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{
				`"uss://test01{asset_type=equity` + BlackBullet + NoSpace,
				`"uss://test01{asset_type=equity` + WhiteBullet + NoSpace,
			},
			v.completeUriAndFile(`"uss://test01{asset_type=e`),
		)
		assert.Equal(t, 0, v.fileMatchCalls)
	})

	t.Run("attr completion with full attr value", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250809,asset_type=equity}"<NS>`,
				`"uss://test01{date=20250809,asset_type=options}"<NS>`,
			},
			v.completeUriAndFile(`"uss://test01{date=20250809,`),
		)
		assert.Equal(t, 0, v.fileMatchCalls)
	})

	t.Run("attr completion with full attr asset_type, date last, same input and output", func(t *testing.T) {
		v := newUriValueTest(t)
		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"uss://test01{asset_type=equity,date=`),
		)
		assert.Equal(t, 0, v.fileMatchCalls)
	})
}

func TestUri_Complete(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		v := newUriValueTest(t)

		// not consider
		v.fileList = []string{
			"file01",
			"file02",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://test01{date=20250219,asset_type=at}"`,
			},
			v.completeUri(`"uss://test01{date=20250219,asset_type=at}"`),
		)

		// check input
		assert.Equal(t, "", v.fileMatchInput)
	})
}

func TestUriAndFile_Complete__With_Zsh(t *testing.T) {
	newTest := func(t *testing.T) *uriValueTest {
		IsZshShellFunc = func() bool {
			return true
		}
		t.Cleanup(func() {
			IsZshShellFunc = isZshShell
		})
		return newUriValueTest(t)
	}

	t.Run("empty", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`"uss://<NS>`,
			},
			v.completeUriAndFile(``),
		)
	})

	t.Run("with uss prefix", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://<NS>`,
			},
			v.completeUriAndFile(`uss`),
		)
	})

	t.Run("with uss prefix and quote", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://<NS>`,
			},
			v.completeUriAndFile(`"uss`),
		)
	})

	t.Run("with empty attribute", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hello{date=<NS>`,
				`uss://hello{asset_type=equity<NS>`,
				`uss://hello{asset_type=options<NS>`,
			},
			v.completeUriAndFile(`"uss://hello{`),
		)
	})

	t.Run("second attribute", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hello{asset_type=equity,date=<NS>`,
			},
			v.completeUriAndFile(`"uss://hello{asset_type=equity,d`),
		)
	})

	t.Run("full", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hello{asset_type=equity,date=20250102}`,
			},
			v.completeUriAndFile(`"uss://hello{asset_type=equity,date=20250102}`),
		)
	})

	t.Run("full, with files", func(t *testing.T) {
		v := newTest(t)
		v.fileList = []string{
			"file01",
			"file02",
		}
		assert.Equal(
			t,
			[]string{
				`uss://hello{asset_type=equity,date=20250102}`,
				`uss://hello{asset_type=equity,date=20250102}/file01`,
				`uss://hello{asset_type=equity,date=20250102}/file02`,
			},
			v.completeUriAndFile(`"uss://hello{asset_type=equity,date=20250102}`),
		)
	})

	t.Run("with dataset name only", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hello`,
			},
			v.completeUriAndFile(`uss://hello`),
		)
	})

	t.Run("with empty attribute, without prefix quote", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hellodate=<NS>`,
				`uss://helloasset_type=equity<NS>`,
				`uss://helloasset_type=options<NS>`,
			},
			v.completeUriAndFile(`uss://hello{`),
		)
	})

	t.Run("with second attr, without prefix quote", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://hellodate=<NS>`,
			},
			v.completeUriAndFile(`uss://hello{asset_type=equity,d`),
		)
	})

	t.Run("with uss prefix and single quote", func(t *testing.T) {
		v := newTest(t)
		assert.Equal(
			t,
			[]string{
				`uss://<NS>`,
			},
			v.completeUriAndFile(`'uss`),
		)
	})
}
