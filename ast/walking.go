package ast

// WalkContext represents the context of a node during a walk. Is given to the WalkFunc.
type WalkContext struct {
	Node *Node
	Path []*Node
}

// WalkFunc is a callback function used in Walk and WalkList.
type WalkFunc func(ctx *WalkContext) error

// Walk performs a pre-order traversal of the AST starting from the given node.
func Walk(n *Node, fn WalkFunc) error {
	ctx := &WalkContext{
		Node: n,
		Path: []*Node{},
	}

	return walk(ctx, fn)
}

func walk(ctx *WalkContext, fn WalkFunc) error {
	if err := fn(ctx); err != nil {
		return err
	}

	for _, child := range ctx.Node.Children {
		childCtx := &WalkContext{
			Node: child,
			Path: append(ctx.Path, ctx.Node),
		}

		if err := walk(childCtx, fn); err != nil {
			return err
		}
	}

	return nil
}

// WalkList performs a pre-order traversal of the AST starting from the given nodes.
func WalkList(nodes []*Node, fn WalkFunc) error {
	for _, n := range nodes {
		if err := Walk(n, fn); err != nil {
			return err
		}
	}

	return nil
}

// WalkFunc is a callback function used in WalkReplace and WalkReplaceList.
// If returning nil, the node being walked is unchanged and won't be replaced.
// If returning []*Node, the node being walked is replaced with the returned nodes.
type WalkReplaceFunc func(ctx *WalkContext) ([]*Node, error)

// WalkReplace performs a post-order traversal of the AST starting from the given node, replacing nodes in-place as returned by WalkReplaceFunc.
func WalkReplace(n *Node, fn WalkReplaceFunc) ([]*Node, error) {
	ctx := &WalkContext{
		Node: n,
		Path: []*Node{},
	}

	return walkReplace(ctx, fn)
}

func walkReplace(ctx *WalkContext, fn WalkReplaceFunc) ([]*Node, error) {
	var newChildren []*Node

	for _, child := range ctx.Node.Children {
		childCtx := &WalkContext{
			Node: child,
			Path: append(ctx.Path, ctx.Node),
		}

		repl, err := walkReplace(childCtx, fn)
		if err != nil {
			return nil, err
		}

		newChildren = append(newChildren, repl...)
	}

	ctx.Node.Children = newChildren

	repl, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	if repl != nil {
		var out []*Node

		for _, r := range repl {
			rCtx := &WalkContext{
				Node: r,
				Path: append(ctx.Path, ctx.Node),
			}

			result, err := walkReplace(rCtx, fn)
			if err != nil {
				return nil, err
			}

			out = append(out, result...)
		}

		return out, nil
	}

	return []*Node{ctx.Node}, nil
}

// WalkReplaceList performs a post-order traversal of the AST starting from the given nodes, replacing nodes in-place as returned by WalkReplaceFunc.
func WalkReplaceList(nodes []*Node, fn WalkReplaceFunc) ([]*Node, error) {
	var out []*Node

	for _, n := range nodes {
		result, err := WalkReplace(n, fn)
		if err != nil {
			return nil, err
		}

		out = append(out, result...)
	}

	return out, nil
}
