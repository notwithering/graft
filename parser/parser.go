package parser

import (
	"errors"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/stack"
	"github.com/notwithering/graft/token"
)

// BuildTree builds an AST from a slice of tokens.
// The blocks map should be used to define which commands are block commands.
func BuildTree(tokens []*token.Token, blocks map[string]bool) ([]*ast.Node, error) {
	root := []*ast.Node{}
	var nodeStack stack.Stack[*ast.Node]

	current := &root

	for _, tok := range tokens {
		switch tok.Kind {
		case token.TokenText:
			*current = append(*current, &ast.Node{
				Kind: ast.NodeText,
				Text: tok.Text,
			})
		case token.TokenOpen:
			nod := &ast.Node{
				Kind:    ast.NodeCommand,
				Command: tok.Command,
				Data:    tok.Data,
				Text:    tok.Text,
			}

			*current = append(*current, nod)

			if blocks[tok.Command] {
				nod.Children = []*ast.Node{}
				nodeStack.Push(nod)
				current = &nod.Children
			}
		case token.TokenClose:
			if nodeStack.Len() == 0 {
				return nil, errors.New("unmatched end")
			}
			nodeStack.Pop()

			if nodeStack.Len() == 0 {
				current = &root
			} else {
				current = &nodeStack.Top().Children
			}
		}
	}

	if nodeStack.Len() != 0 {
		return nil, errors.New("unclosed blocks")
	}

	return root, nil
}
