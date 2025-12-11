package settings

import (
	"strings"

	"golang.org/x/xerrors"
)

const (
	BehaviorBefore = "before"
	BehaviorAfter  = "after"
)

type AtcSettings struct {
	Path     string `yaml:"path"`
	Behavior string `yaml:"behavior"`
	Template string `yaml:"template"`
	Branch   string `yaml:"branch"`
	RegexStr string `yaml:"regexstr"`
}

func validateSettings(settings *AtcSettings) error {
	//check Behavior:
	if strings.ToLower(settings.Behavior) != BehaviorAfter && strings.ToLower(settings.Behavior) != BehaviorBefore {
		return xerrors.Errorf("invalid behavior")
	}
	//check Template:
	if !strings.Contains(settings.Template, `{{.Version}}`) {
		return xerrors.Errorf("template does not contain {{.Version}}")
	}
	//check Path:
	pathPrefix := "/"

	if settings.Path == "" {
		return nil
	}
	if strings.HasPrefix(settings.Path, pathPrefix) {
		return xerrors.Errorf("invalid path prefix")
	}
	if strings.Contains(settings.Path, "//") {
		return xerrors.Errorf("invalid path format")
	}
	return nil
}
