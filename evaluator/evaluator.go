package evaluator

import (
	"monkey/ast"
	"monkey/obj"
)

func Eval(node ast.Node) obj.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expr)

	// Expressions
	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(statements []ast.Statement) obj.Object {
	var ret obj.Object

	for _, s := range statements {
		ret = Eval(s)
	}
	return ret
}
