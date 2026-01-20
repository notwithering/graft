package ast

type WalkContext struct {
	Node *Node
	Path []*Node
}

type WalkFunc func(ctx *WalkContext) error

// pre order
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

// pre order
func WalkList(nodes []*Node, fn WalkFunc) error {
	for _, n := range nodes {
		if err := Walk(n, fn); err != nil {
			return err
		}
	}

	return nil
}

type WalkReplaceFunc func(ctx *WalkContext) ([]*Node, error)

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
