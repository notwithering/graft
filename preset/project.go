package preset

import (
	"errors"

	"github.com/notwithering/graft/ast"
)

type Project struct {
	Config        ProjectConfig
	Sources       map[string]*Source
	NodeSourceMap map[*ast.Node]*Source
}

func NewProject(projectConfig ProjectConfig) *Project {
	return &Project{
		Config:        projectConfig,
		Sources:       make(map[string]*Source),
		NodeSourceMap: make(map[*ast.Node]*Source),
	}
}

type ProjectConfig struct {
	Root     string
	Commands map[string]*CommandSpec
}

// GetSpecName iterates over the projects commands until it finds the target CommandSpec, and
// returns its defined name.
func (proj *Project) GetSpecName(target *CommandSpec) (string, error) {
	for name, spec := range proj.Config.Commands {
		if spec == target {
			return name, nil
		}
	}

	return "", errors.New("target not found in project commands")
}
