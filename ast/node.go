package ast

type Node interface {
	Literal() string
	String() string
}
