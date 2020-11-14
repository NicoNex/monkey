package ast

type Expression interface {
	Node
	ENode()
}
