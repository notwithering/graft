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

func (proj *Project) GetSpecName(target *CommandSpec) (string, error) {
	for name, spec := range proj.Config.Commands {
		if spec == target {
			return name, nil
		}
	}

	return "", errors.New("target not found in project commands")
}
