package settings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSettingsForErrors(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		behavior string
		template string
		branch   string
		regexstr string
		wantErr  string
	}{
		{
			name:     "invalid path prefix",
			path:     "/contents/pom.xml",
			behavior: "before",
			template: "v{{.Version}}",
			branch:   "",
			regexstr: "",
			wantErr:  "invalid path prefix",
		},
		{
			name:     "invalid path format",
			path:     "contents//asd.txt",
			behavior: "before",
			template: "v{{.Version}}",
			branch:   "",
			regexstr: "",
			wantErr:  "invalid path format",
		},
		{
			name:     "valid path",
			path:     "contents/asd.txt",
			behavior: "before",
			template: "v{{.Version}}",
			branch:   "",
			regexstr: "",
		},
		{
			name:     "invalid behavior",
			path:     "contents/pom.xml",
			behavior: "bef",
			template: "v{{.Version}}",
			branch:   "",
			regexstr: "",
			wantErr:  "invalid behavior",
		},
		{
			name:     "template missing version placeholder",
			path:     "package.json",
			behavior: "after",
			template: "{.version}",
			branch:   "",
			regexstr: "",
			wantErr:  "template does not contain {{.Version}}",
		},
		{
			name:     "valid settings before",
			path:     "contents/pom.xml",
			behavior: "before",
			template: "v{{.Version}}V",
			branch:   "testbranch",
			regexstr: "",
		},
		{
			name:     "valid settings after",
			path:     "package.json",
			behavior: "after",
			template: "{{.Version}}-release",
			branch:   "main",
			regexstr: "",
		},
		{
			name:     "empty path",
			path:     "",
			behavior: "before",
			template: "v{{.Version}}",
			branch:   "",
			regexstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := &AtcSettings{
				Path:     tt.path,
				Behavior: tt.behavior,
				Template: tt.template,
				Branch:   tt.branch,
				RegexStr: tt.regexstr,
			}
			err := validateSettings(settings)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
