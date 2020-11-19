package evaluator

import (
	"fmt"
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

func evalProgram(statements []ast.Statement) obj.Object {
	var ret obj.Object

	for _, s := range statements {
		ret = Eval(s)

		switch result := ret.(type) {
		case *obj.ReturnValue:
			return result.Value
		case *obj.Error:
			return result
		}
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
		return newError("unknown operator: -%s", right.Type().String())
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
		return newError("unknown operator %s%s", op, right.Type().String())
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

	case left.Type() != right.Type():
		lt := left.Type().String()
		rt := right.Type().String()
		return newError("type mismatch: %s %s %s", lt, op, rt)

	default:
		lt := left.Type().String()
		rt := right.Type().String()
		return newError("unknown operator: %s %s %s", lt, op, rt)
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
		lt := left.Type().String()
		rt := right.Type().String()
		return newError("unknown operator: %s %s %s", lt, op, rt)
	}
}

func evalIfExpr(ie *ast.IfExpression) obj.Object {
	cond := Eval(ie.Condition)

	if isError(cond) {
		return cond
	}

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

func evalBlockStatement(block *ast.BlockStatement) obj.Object {
	var res obj.Object

	for _, s := range block.Statements {
		res = Eval(s)

		if res != nil {
			rt := res.Type()
			if rt == obj.RETURN || rt == obj.ERROR {
				return res
			}
		}
	}
	return res
}

func newError(format string, a ...interface{}) *obj.Error {
	return &obj.Error{Msg: fmt.Sprintf(format, a...)}
}

func isError(o obj.Object) bool {
	if o != nil {
		return o.Type() == obj.ERROR
	}
	return false
}

func Eval(node ast.Node) obj.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expr)

	// Expressions
	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}

	case *ast.Boolean:
		return btoo(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpr(node)

	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &obj.ReturnValue{Value: val}
	}

	return nil
}
