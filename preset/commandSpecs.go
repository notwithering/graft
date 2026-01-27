package preset

import (
	"github.com/notwithering/graft/ast"
)

var DefaultCommands = map[string]*CommandSpec{
	"include": IncludeCommand,

	"extend": ExtendCommand,
	"define": DefineCommand,
	"block":  BlockCommand,
}

var IncludeCommand = &CommandSpec{
	Args: map[string]ArgType{
		"src": ArgTypeSourcePtr,
	},
	Block: false,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		target := ctx.Args["src"].(*Source)
		return target.Tree, nil
	},
}

var ExtendCommand = &CommandSpec{
	Args: map[string]ArgType{
		"src": ArgTypeSourcePtr,
	},
	Block: true,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		target := ctx.Args["src"].(*Source)

		definitions := make(map[string]*ast.Node)

		ast.WalkList(ctx.Node.Children, func(ctx *ast.WalkContext) error {
			if ctx.Node.Kind != ast.NodeCommand || ctx.Node.Command != "define" {
				return nil
			}

			definitions[ctx.Node.Data.(map[string]any)["name"].(string)] = ctx.Node
			return nil
		})

		tree, err := ast.WalkReplaceList(target.Tree, func(ctx *ast.WalkContext) ([]*ast.Node, error) {
			if ctx.Node.Kind != ast.NodeCommand || ctx.Node.Command != "block" {
				return nil, nil
			}

			def, ok := definitions[ctx.Node.Data.(map[string]any)["name"].(string)]
			if !ok {
				return nil, nil
			}

			return def.Children, nil
		})
		if err != nil {
			return nil, err
		}

		return tree, nil
	},
}

var DefineCommand = &CommandSpec{
	Args: map[string]ArgType{
		"name": ArgTypeString,
	},
	Block: true,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		return nil, nil
	},
}

var BlockCommand = &CommandSpec{
	Args: map[string]ArgType{
		"name": ArgTypeString,
	},
	Block: false,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		return nil, nil
	},
}
