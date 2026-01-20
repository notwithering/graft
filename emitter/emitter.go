package emitter

import (
	"strings"

	"github.com/notwithering/graft/ast"
)

func Emit(tree []*ast.Node) string {
	var out strings.Builder
	emitNodes(&out, tree)
	return out.String()
}

func emitNode(sb *strings.Builder, n *ast.Node) {
	sb.WriteString(n.Text)
	emitNodes(sb, n.Children)
}

func emitNodes(sb *strings.Builder, tree []*ast.Node) {
	for _, n := range tree {
		emitNode(sb, n)
	}
}
