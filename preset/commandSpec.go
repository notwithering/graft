package preset

import (
	"github.com/notwithering/graft/ast"
)

type CommandSpec struct {
	Args    []ArgType
	Block   bool
	Rewrite RewriteFunc
}

type ArgType uint8

const (
	ArgTypeString ArgType = iota
	ArgTypeSourcePtr
)

type RewriteFunc func(ctx *Context) ([]*ast.Node, error)

type Context struct {
	Project *Project
	Source  *Source
	Node    *ast.Node
	Args    []any
}
