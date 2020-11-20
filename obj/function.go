package obj

import (
	"fmt"
	"monkey/ast"
	"strings"
)

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var params []string

	for _, p := range f.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body)
}
