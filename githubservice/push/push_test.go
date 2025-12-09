package push

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMadeСaptionToTemplate(t *testing.T) {
	var tests = []struct {
		template string
		version  string
		result   string
	}{
		{`v{{.Version}}`, `1.0`, `v1.0`},
		{`vNN{{.Version}}`, `1.0`, `vNN1.0`},
		{`v_{{.Version}}`, `1.0`, `v_1.0`},
		{`v{{.Version}}`, `1.0-relise`, `v1.0-relise`},
		{`{{.Version}}`, `1.0`, `1.0`},
		{`Time hour now: {{Time.Hour}}, {{.Version}}`, `1.0`, "Time hour now: " + strconv.Itoa(time.Now().Hour()) + ", 1.0"},
		{``, `1.0`, ``},
	}
	for _, test := range tests {
		result, _ := renderTagNameTemplate(test.template, test.version)
		if result != test.result {
			t.Errorf("template: %q, version: %q\nwant: %q, got: %q", test.template, test.version, test.result, result)
		}
	}
}

func TestMadeСaptionToTemplateError(t *testing.T) {
	var tests = []struct {
		template  string
		version   string
		errString string
	}{
		{`v{{.Versio}}`, `1.0`, `template: template tagContent:1:3: executing "template tagContent" at <.Versio>: can't evaluate field Versio in type push.TagContent`},
	}
	for _, test := range tests {
		_, err := renderTagNameTemplate(test.template, test.version)
		if fmt.Sprint(err) != test.errString {
			t.Errorf("template: %q, version: %q\nerr want: %v, err got: %v", test.template, test.version, test.errString, err)
		}
	}
}
