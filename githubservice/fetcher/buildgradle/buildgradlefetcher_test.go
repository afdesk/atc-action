package buildgradle

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalBuildGradle(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		version string
		wantErr string
	}{
		{
			name:    "happy path",
			path:    "testdata/happy",
			version: "1.7.1",
		},
		{
			name:    "happy path with release",
			path:    "testdata/happy-release",
			version: "1.7.1-release",
		},
		{
			name:    "happy path with comment",
			path:    "testdata/happy-comment",
			version: "1.7.1",
		},
		{
			name:    "sad path",
			path:    "testdata/sad",
			wantErr: "empty number version",
		},
	}
	for _, tt := range tests {
		gradle := &BuildGradle{}
		content, err := os.ReadFile(tt.path)
		err = unmarshalBuildGradle(content, gradle)
		if tt.wantErr != "" {
			require.EqualError(t, err, tt.wantErr)
			return
		}
		require.NoError(t, err)
		require.Equal(t, tt.version, gradle.Version)

	}
}
