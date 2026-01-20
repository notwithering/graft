package main

import (
	"github.com/notwithering/graft/preset"
	"github.com/notwithering/graft/syntax"
)

var (
	root     = "./src"
	syntaxes = map[string]*syntax.Syntax{
		"html": syntax.SGMLTagSyntax,
	}
	commands = preset.DefaultCommands
	dest     = "./dist"
)

func main() {
	proj := preset.NewProject(preset.ProjectConfig{
		Root: root,
	})

	if err := proj.Assemble(syntaxes, commands); err != nil {
		panic(err)
	}

	if err := proj.Resolve(commands); err != nil {
		panic(err)
	}

	if err := proj.Write(dest); err != nil {
		panic(err)
	}
}
