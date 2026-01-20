package preset

import (
	"strings"

	"github.com/notwithering/graft/ast"
)

func (proj *Project) showCycle(ctx *ast.WalkContext) string {
	index := -1

	for i, n := range ctx.Path {
		if n == ctx.Node {
			index = i
			break
		}
	}

	if index == -1 {
		var chain []string

		for _, n := range append(ctx.Path, ctx.Node) {
			chain = append(chain, n.String())
		}

		return strings.Join(chain, " -> ")
	}

	cycle := ctx.Path[index:]

	var out strings.Builder

	for i, node := range cycle {
		if len(cycle) == 1 {
			out.WriteString("↪")
		} else if i == 0 {
			out.WriteString("┌→┌─")
		} else if i != len(cycle)-1 {
			out.WriteString("│ └→")
		} else {
			out.WriteString("└─└→")
		}

		src, ok := proj.NodeSourceMap[node]
		if ok {
			out.WriteString(src.LocalPath)
			out.WriteString(": ")
		}

		out.WriteString(node.Text)
		out.WriteByte('\n')
	}

	return out.String()
}
