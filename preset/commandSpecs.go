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
	Args:  []ArgType{ArgTypeSourcePtr},
	Block: false,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		target := ctx.Args[0].(*Source)
		return target.Tree, nil
	},
}

var ExtendCommand = &CommandSpec{
	Args:  []ArgType{ArgTypeSourcePtr},
	Block: true,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		target := ctx.Args[0].(*Source)

		definitions := make(map[string]*ast.Node)

		ast.WalkList(ctx.Node.Children, func(ctx *ast.WalkContext) error {
			if ctx.Node.Kind != ast.NodeCommand || ctx.Node.Args[0] != "define" || len(ctx.Node.Args) != 2 {
				return nil
			}

			definitions[ctx.Node.Args[1]] = ctx.Node
			return nil
		})

		tree, err := ast.WalkReplaceList(target.Tree, func(ctx *ast.WalkContext) ([]*ast.Node, error) {
			if ctx.Node.Kind != ast.NodeCommand || ctx.Node.Args[0] != "block" || len(ctx.Node.Args) != 2 {
				return nil, nil
			}

			def, ok := definitions[ctx.Node.Args[1]]
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
	Args:  []ArgType{ArgTypeString},
	Block: true,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		return nil, nil
	},
}

var BlockCommand = &CommandSpec{
	Args:  []ArgType{ArgTypeString},
	Block: false,
	Rewrite: func(ctx *Context) ([]*ast.Node, error) {
		return nil, nil
	},
}
