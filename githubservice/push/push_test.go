package push

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMadeСaptionToTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		version  string
		result   string
		wantErr  bool
	}{
		{
			name:     "simple version",
			template: `v{{.Version}}`,
			version:  `1.0`,
			result:   `v1.0`,
		},
		{
			name:     "version with prefix",
			template: `vNN{{.Version}}`,
			version:  `1.0`,
			result:   `vNN1.0`,
		},
		{
			name:     "version with underscore",
			template: `v_{{.Version}}`,
			version:  `1.0`,
			result:   `v_1.0`,
		},
		{
			name:     "version with release",
			template: `v{{.Version}}`,
			version:  `1.0-relise`,
			result:   `v1.0-relise`,
		},
		{
			name:     "only version placeholder",
			template: `{{.Version}}`,
			version:  `1.0`,
			result:   `1.0`,
		},
		{
			name:     "version with time placeholder",
			template: `Time hour now: {{Time.Hour}}, {{.Version}}`,
			version:  `1.0`,
			result:   "Time hour now: " + strconv.Itoa(time.Now().Hour()) + ", 1.0"},
		{
			name:     "empty template",
			template: ``,
			version:  `1.0`,
			result:   ``,
		},
		{
			name:     "invalid field",
			template: `v{{.Versio}}`,
			version:  `1.0`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderTagNameTemplate(tt.template, tt.version)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		})
	}
}
