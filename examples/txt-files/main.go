package main

import (
	"github.com/notwithering/graft/preset"
	"github.com/notwithering/graft/token"
)

var (
	root     = "./src"
	syntaxes = map[string]*token.Syntax{
		"txt": token.DoubleBraceSyntax,
	}
	commands = preset.DefaultCommands
	dest     = "./dist"
)

func main() {
	proj := preset.NewProject(preset.ProjectConfig{
		Root:     root,
		Commands: commands,
	})

	if err := proj.Assemble(syntaxes); err != nil {
		panic(err)
	}

	if err := proj.Resolve(); err != nil {
		panic(err)
	}

	if err := proj.Write(dest); err != nil {
		panic(err)
	}
}
