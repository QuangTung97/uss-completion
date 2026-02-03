package completion

import (
	"fmt"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func TestUriAndFile_Complete__Basic(t *testing.T) {
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

	matchDatasetInputs  []string
	matchDatasetOutputs []string

	listVersionNames []string
}

func newUriValueTest(t *testing.T) *uriValueTest {
	v := &uriValueTest{}

	// stub get dataset names
	GetMatchDatasetNamesFunc = func(match string) []string {
		v.matchDatasetInputs = append(v.matchDatasetInputs, match)
		return v.matchDatasetOutputs
	}
	t.Cleanup(func() {
		GetMatchDatasetNamesFunc = getMatchDatasetNamesTest
	})

	// stub get files
	globalListFilesByPatternFunc = func(dir string, match string) []string {
		v.fileMatchDir = dir
		v.fileMatchInput = match
		v.fileMatchCalls++
		return v.fileList
	}
	t.Cleanup(func() {
		globalListFilesByPatternFunc = listFilesByPattern
	})

	// stub get uri path
	GetUriDiskPathFunc = func(_ string) string {
		return "uss_storage"
	}
	t.Cleanup(func() {
		GetUriDiskPathFunc = getUriDiskPathTest
	})

	// stub get versions
	GetAllVersionsFunc = func(dsName string) VersionList {
		v.listVersionNames = append(v.listVersionNames, dsName)
		return getAllVersionsTest(dsName)
	}
	t.Cleanup(func() {
		GetAllVersionsFunc = getAllVersionsTest
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
		// check list versions input
		assert.Equal(t, []string{"test01"}, v.listVersionNames)
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

	t.Run("disable attr completion", func(t *testing.T) {
		v := newUriValueTest(t)

		GetAllVersionsFunc = func(_ string) VersionList {
			return VersionList{
				DisableCompletion: true,
			}
		}

		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"uss://test01{`),
		)
		assert.Equal(t, 0, v.fileMatchCalls)
	})
}

func TestUri_Complete(t *testing.T) {
	t.Run("match exact", func(t *testing.T) {
		v := newUriValueTest(t)

		// not consider
		v.fileList = []string{
			"file01",
			"file02",
		}

		assert.Equal(
			t,
			[]string{},
			v.completeUri(`"uss://test01{date=20250219,asset_type=at}"`),
		)

		// check input
		assert.Equal(t, "", v.fileMatchInput)
	})

	t.Run("dataset name completion", func(t *testing.T) {
		v := newUriValueTest(t)

		v.matchDatasetOutputs = []string{
			"dataset01",
			"dataset02",
		}

		assert.Equal(
			t,
			[]string{
				`"uss://dataset01{<NS>`,
				`"uss://dataset02{<NS>`,
			},
			v.completeUriAndFile(`"uss://data`),
		)

		assert.Equal(t, []string{"data"}, v.matchDatasetInputs)
	})

	t.Run("dataset name completion, reach limit 20", func(t *testing.T) {
		v := newUriValueTest(t)

		for i := range 20 {
			v.matchDatasetOutputs = append(v.matchDatasetOutputs, fmt.Sprintf("dataset%02d", i+1))
		}

		var expected []string
		for i := range 20 {
			expected = append(expected, fmt.Sprintf(`"uss://dataset%02d{<NS>`, i+1))
		}
		expected = append(expected, `"uss://<NS>`)

		assert.Equal(
			t,
			expected,
			v.completeUriAndFile(`"uss://data`),
		)

		assert.Equal(t, []string{"data"}, v.matchDatasetInputs)
	})

	t.Run("prefix is not uss", func(t *testing.T) {
		v := newUriValueTest(t)

		assert.Equal(
			t,
			[]string{},
			v.completeUriAndFile(`"ss://data`),
		)
	})

	t.Run("all attrs, complete bracket", func(t *testing.T) {
		v := newUriValueTest(t)

		assert.Equal(
			t,
			[]string{
				`"uss://dataset01{date=20260109,asset_type=equity}"` + BlackBullet + NoSpace,
				`"uss://dataset01{date=20260109,asset_type=equity}"` + WhiteBullet + NoSpace,
			},
			v.completeUriAndFile(`"uss://dataset01{date=20260109,asset_type=equity`),
		)
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
				`uss://hello{asset_type=equity,date=20250102}<NS>`,
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
			[]string{},
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

	t.Run("with get uri path return null dir", func(t *testing.T) {
		v := newTest(t)

		GetUriDiskPathFunc = func(_ string) string {
			return NullDir
		}

		assert.Equal(
			t,
			[]string{
				`uss://hello{date=20250912}<NS>`,
			},
			v.completeUriAndFile(`'uss://hello{date=20250912}`),
		)

		assert.Equal(t, 0, v.fileMatchCalls)
		assert.Equal(t, "", v.fileMatchDir)
	})
}
