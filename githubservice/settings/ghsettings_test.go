package settings

import (
	"fmt"
	"testing"
)

func TestCheckSettingsForErrors(t *testing.T) {
	var tests = []struct {
		path             string
		behavior         string
		template         string
		branch           string
		regexstr         string
		expectedErrorStr string
	}{
		{"/contents/pom.xml", "", "", "", "", `error config file .atc.yaml; path has prefix "/"`},
		{"contents//asd.txt", "", "", "", "", `error config file .atc.yaml; path has "//"`},
		{"contents/asd.txt", "", "", "", "", fmt.Sprint(nil)},
		{"contents/pom.xml/", "bef", "", "", "", `error config file .atc.yaml: behavior doesn't contain "before" or "after"`},
		{"package.json", "after", "{.version}", "", "", `error config file .atc.yaml: template doesn't contain "{{.Version}}"`},
		{"pubspec.yaml", "before", ".vers", "", "", `error config file .atc.yaml: template doesn't contain "{{.Version}}"`},
		{"contents/pom.xml", "before", "v{{.Version}}V", "testbranch", "", fmt.Sprint(nil)},
	}

	for _, test := range tests {
		settings := &AtcSettings{test.path, test.behavior, test.template, test.branch, test.regexstr}
		err := validateSettings(settings)
		if fmt.Sprint(err) != test.expectedErrorStr {
			t.Errorf("no takes error settings:%s\nexpected: %s, got: %s", settings, test.expectedErrorStr, err)
		}
	}
}
