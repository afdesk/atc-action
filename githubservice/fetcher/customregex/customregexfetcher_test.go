package customregex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalUserConfig(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		regexstr string
		version  string
		wantErr  string
	}{
		{
			name:     "happy path",
			path:     "testdata/happy",
			regexstr: "vers: \"(.+?)\"",
			version:  "1.1",
		},
		{
			name:     "happy with release",
			path:     "testdata/happy-release",
			regexstr: "vers: \"(.+?)\"",
			version:  "1.1-release",
		},
		{
			name:     "invalid regex syntax",
			path:     "",
			regexstr: "vers: (.{}",
			wantErr:  "error parsing regexp: missing closing ): `vers: (.{}`",
		},
		{
			name:     "empty version",
			path:     "testdata/happy",
			regexstr: "versious: 1",
			wantErr:  "empty number version",
		},
		{
			name:     "empty regex group",
			path:     "",
			regexstr: "",
			wantErr:  "regexStr don't have group",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customRegexConf := &Config{}
			content, err := os.ReadFile(tt.path)
			err = unmarshalCustomRegexConfig(content, tt.regexstr, customRegexConf)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.version, customRegexConf.Version)
		})
	}
}
