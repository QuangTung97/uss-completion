package completion

type VersionAttr = map[string]string

type VersionList struct {
	Versions []VersionAttr
}

func getAllVersionsTest() VersionList {
	return VersionList{
		Versions: []VersionAttr{
			{"asset_type": "equity"},
			{"asset_type": "options"},
		},
	}
}

var GetAllVersionsFunc = getAllVersionsTest

func getUriDiskPathTest(string) string {
	return ""
}

var GetUriDiskPathFunc = getUriDiskPathTest
