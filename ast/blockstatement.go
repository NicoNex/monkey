package ast

import (
	"bytes"
	"monkey/token"
)

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

func (b *BlockStatement) SNode() {}

func (b *BlockStatement) Literal() string {
	return b.Token.Lit
}

func (b *BlockStatement) String() string {
	var buf bytes.Buffer

	for _, s := range b.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
}
