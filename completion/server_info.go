package completion

type VersionAttr = map[string]string

type VersionList struct {
	Versions          []VersionAttr
	DisableCompletion bool
}

func getAllVersionsTest(_ string) VersionList {
	return VersionList{
		Versions: []VersionAttr{
			{"asset_type": "equity"},
			{"asset_type": "options"},
		},
	}
}

var GetAllVersionsFunc = getAllVersionsTest

func getUriDiskPathTest(_ string) string {
	return "uss_storage"
}

var GetUriDiskPathFunc = getUriDiskPathTest

func getMatchDatasetNamesTest(_ string) []string {
	return []string{
		"hello01",
		"tung",
		"sample",
	}
}

var GetMatchDatasetNamesFunc = getMatchDatasetNamesTest
