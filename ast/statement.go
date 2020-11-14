package ast

type Statement interface {
	Node
	SNode()
}
