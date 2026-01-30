package preset

import (
	"github.com/notwithering/graft/ast"
)

var DefaultCommands = map[string]*CommandSpec{
	"include": IncludeCommand,

	"extend": ExtendCommand,
	"define": DefineCommand,
	"block":  BlockCommand,

	"delete": DeleteCommand,
}

var IncludeCommand = &CommandSpec{
	Args: map[string]ArgType{
		"src": ArgTypeSourcePtr,
	},
	Block: false,
	Rewrite: func(ctx *CommandContext) ([]*ast.Node, error) {
		target := ctx.Args["src"].(*Source)
		return target.Tree, nil
	},
}

var ExtendCommand = &CommandSpec{
	Args: map[string]ArgType{
		"src": ArgTypeSourcePtr,
	},
	Block: true,
	Rewrite: func(cmdCtx *CommandContext) ([]*ast.Node, error) {
		target := cmdCtx.Args["src"].(*Source)

		definitions := make(map[string]*ast.Node)

		ast.WalkList(cmdCtx.Node.Children, func(walkCtx *ast.WalkContext) error {
			if walkCtx.Node.Kind != ast.NodeCommand || walkCtx.Node.Command != "define" {
				return nil
			}

			newCtx := cmdCtx.Clone()
			newCtx.Node = walkCtx.Node

			args, err := newCtx.ParseArgTypes(DefineCommand.Args)
			if err != nil {
				return err
			}

			definitions[args["name"].(string)] = walkCtx.Node
			return nil
		})

		tree, err := ast.WalkReplaceList(target.Tree, func(walkCtx *ast.WalkContext) ([]*ast.Node, error) {
			if walkCtx.Node.Kind != ast.NodeCommand || walkCtx.Node.Command != "block" {
				return nil, nil
			}

			newCtx := cmdCtx.Clone()
			newCtx.Node = walkCtx.Node

			args, err := newCtx.ParseArgTypes(BlockCommand.Args)
			if err != nil {
				return nil, err
			}

			def, ok := definitions[args["name"].(string)]
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
	Rewrite: func(ctx *CommandContext) ([]*ast.Node, error) {
		return nil, nil
	},
}

var BlockCommand = &CommandSpec{
	Args: map[string]ArgType{
		"name": ArgTypeString,
	},
	Block: false,
	Rewrite: func(ctx *CommandContext) ([]*ast.Node, error) {
		return nil, nil
	},
}

var DeleteCommand = &CommandSpec{
	Args:  map[string]ArgType{},
	Block: true,
	Rewrite: func(ctx *CommandContext) ([]*ast.Node, error) {
		return []*ast.Node{}, nil
	},
}
