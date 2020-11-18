package evaluator

import (
	"monkey/ast"
	"monkey/obj"
)

var (
	NULL  = &obj.Null{}
	TRUE  = &obj.Boolean{Value: true}
	FALSE = &obj.Boolean{Value: false}
)

// Returns the object representation of the boolean primitive b.
func btoo(b bool) *obj.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalStatements(statements []ast.Statement) obj.Object {
	var ret obj.Object

	for _, s := range statements {
		ret = Eval(s)
	}
	return ret
}

func evalBangOpExpr(right obj.Object) obj.Object {
	switch right {
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalPrefixMinusOpExpr(right obj.Object) obj.Object {
	if right.Type() != obj.INT {
		return NULL
	}

	val := right.(*obj.Integer).Value
	return &obj.Integer{Value: -val}
}

func evalPrefixExpr(op string, right obj.Object) obj.Object {
	switch op {
	case "!":
		return evalBangOpExpr(right)

	case "-":
		return evalPrefixMinusOpExpr(right)

	default:
		return NULL
	}
}

func evalInfixExpr(op string, left, right obj.Object) obj.Object {
	switch {
	case left.Type() == obj.INT && right.Type() == obj.INT:
		return evalIntInfixExpr(op, left, right)
	case op == "==":
		return btoo(left == right)
	case op == "!=":
		return btoo(left != right)
	default:
		return NULL
	}
}

func evalIntInfixExpr(op string, left, right obj.Object) obj.Object {
	var l = left.(*obj.Integer).Value
	var r = right.(*obj.Integer).Value

	switch op {
	case "+":
		return &obj.Integer{Value: l + r}
	case "-":
		return &obj.Integer{Value: l - r}
	case "*":
		return &obj.Integer{Value: l * r}
	case "/":
		return &obj.Integer{Value: l / r}
	case "==":
		return btoo(l == r)
	case "!=":
		return btoo(l != r)
	case "<":
		return btoo(l < r)
	case ">":
		return btoo(l > r)
	case "<=":
		return btoo(l <= r)
	case ">=":
		return btoo(l >= r)
	default:
		return NULL
	}
}

func evalIfExpr(ie *ast.IfExpression) obj.Object {
	cond := Eval(ie.Condition)

	if isTruthy(cond) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}
	return NULL
}

func isTruthy(cond obj.Object) bool {
	return cond != NULL && cond != FALSE
}

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

	case *ast.Boolean:
		return btoo(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.IfExpression:
		return evalIfExpr(node)
	}

	return nil
}
