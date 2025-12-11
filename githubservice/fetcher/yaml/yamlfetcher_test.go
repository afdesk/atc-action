package yaml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalYaml(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		version string
		wantErr string
	}{
		{
			name:    "happy path",
			path:    "testdata/happy",
			version: "1.2.0",
		},
		{
			name:    "happy with release",
			path:    "testdata/happy-release",
			version: "1.2.0-release",
		},
		{
			name:    "invalid format",
			path:    "testdata/sad",
			wantErr: "yaml: line 4: could not find expected ':'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubspecyaml := &Yaml{}
			content, err := os.ReadFile(tt.path)
			err = UnmarshalYaml(content, pubspecyaml)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.version, pubspecyaml.Version)
		})
	}
}
