package ast

import "fmt"

type Node struct {
	Kind NodeKind

	Command string
	Data    any
	Text    string

	Children []*Node

	Resolved  bool
	Resolving bool
}

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
		Kind:      n.Kind,
		Command:   n.Command,
		Data:      n.Data,
		Text:      n.Text,
		Children:  append([]*Node{}, n.Children...),
		Resolved:  n.Resolved,
		Resolving: n.Resolving,
	}
}
