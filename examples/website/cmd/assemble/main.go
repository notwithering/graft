package main

import (
	"github.com/notwithering/graft/preset"
	"github.com/notwithering/graft/token"
)

var (
	root     = "./src"
	syntaxes = map[string]*token.Syntax{
		"html": token.HTMLTagSyntax,
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
