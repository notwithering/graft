package preset

import (
	"fmt"
	"reflect"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/pathutil"
)

type CommandSpec struct {
	Args    map[string]ArgType
	Block   bool
	Rewrite RewriteFunc
}

type ArgType uint8

const (
	ArgTypeString ArgType = iota
	ArgTypeSourcePtr
)

type RewriteFunc func(ctx *CommandContext) ([]*ast.Node, error)

type CommandContext struct {
	Project *Project
	Source  *Source
	Node    *ast.Node
	Args    map[string]any
}

func (ctx CommandContext) Clone() *CommandContext {
	return &CommandContext{
		Project: ctx.Project,
		Source:  ctx.Source,
		Node:    ctx.Node,
		Args:    ctx.Args,
	}
}

func (ctx *CommandContext) ParseArgTypes(argTypes map[string]ArgType) (map[string]any, error) {
	args := make(map[string]any)

	for key, arg := range ctx.Node.Data.(map[string]any) {
		argType, ok := argTypes[key]
		if !ok {
			continue
		}

		switch argType {
		case ArgTypeString:
			argStr, ok := arg.(string)
			if !ok {
				return nil, fmt.Errorf("incompatible type: %s", reflect.TypeOf(arg))
			}

			args[key] = argStr
		case ArgTypeSourcePtr:
			nsrc, ok := ctx.Project.NodeSourceMap[ctx.Node]
			if !ok {
				return nil, fmt.Errorf("source mapping not found: %s", ctx.Node)
			}

			argStr, ok := arg.(string)
			if !ok {
				return nil, fmt.Errorf("incompatible type: %s", reflect.TypeOf(arg))
			}
			targetPath := pathutil.TargetPath(nsrc.LocalPath, argStr)

			targetSource, ok := ctx.Project.Sources[targetPath]
			if !ok {
				return nil, fmt.Errorf("target not found: %s", targetPath)
			}

			args[key] = targetSource
		}
	}

	return args, nil
}
