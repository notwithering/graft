package ast

import "fmt"

// Node is a single AST element.
// Interpretation depends on Kind.
type Node struct {
	Kind NodeKind

	Command string
	Data    any
	Text    string // Text includes full match by syntax.Syntax.OpenClose.

	Children []*Node
}

// NodeKind represents the type of an AST node. Determines how it should be interpreted.
type NodeKind uint8

const (
	NodeText NodeKind = iota
	NodeCommand
)

func (n Node) String() string {
	switch n.Kind {
	case NodeText:
		return fmt.Sprintf("NodeText(%q)", n.Text)
	case NodeCommand:
		return fmt.Sprintf("NodeCommand(%q)", n.Text)
	default:
		return fmt.Sprintf("Node(%d)", n.Kind)
	}
}

func (n Node) Clone() *Node {
	return &Node{
		Kind:     n.Kind,
		Command:  n.Command,
		Data:     n.Data,
		Text:     n.Text,
		Children: append([]*Node{}, n.Children...),
	}
}
