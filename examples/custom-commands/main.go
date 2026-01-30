package main

import (
	"time"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/preset"
	"github.com/notwithering/graft/token"
)

var (
	root     = "./src"
	syntaxes = map[string]*token.Syntax{
		"txt": token.DoubleBraceSyntax,
	}
	commands = map[string]*preset.CommandSpec{
		"extend": preset.ExtendCommand,
		"define": preset.DefineCommand,
		"block":  preset.BlockCommand,
		"last-built": {
			Args:  map[string]preset.ArgType{},
			Block: false,
			Rewrite: func(ctx *preset.CommandContext) ([]*ast.Node, error) {
				return []*ast.Node{{
					Kind: ast.NodeText,
					Text: time.Now().Format("2006-01-02 15:04:05"),
				}}, nil
			},
		},
	}
	dest = "./dist"
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
