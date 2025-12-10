package pomxml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalPomXml(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		version string
		wantErr string
	}{
		{
			name:    "happy path",
			path:    "testdata/happy",
			version: "4.2.0",
		},
		{
			name:    "version with release",
			path:    "testdata/happy-release",
			version: "4.2.0-release",
		},
		{
			name:    "xml declaration only",
			path:    "testdata/sad",
			wantErr: "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pomxml := &PomXml{}
			content, err := os.ReadFile(tt.path)
			err = unmarshalPomXml(content, pomxml)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.version, pomxml.Version)
		})
	}
}
