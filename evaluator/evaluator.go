package evaluator

import (
	"fmt"
	"github.com/NicoNex/monkey/ast"
	"github.com/NicoNex/monkey/obj"
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

func evalProgram(statements []ast.Statement, env *obj.Env) obj.Object {
	var ret obj.Object

	for _, s := range statements {
		ret = Eval(s, env)

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

	case left.Type() == obj.STRING && right.Type() == obj.STRING:
		return evalStrInfixExpr(op, left, right)

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

func evalStrInfixExpr(op string, left, right obj.Object) obj.Object {
	var l = left.(*obj.String).Value
	var r = right.(*obj.String).Value

	switch op {

	case "+":
		return &obj.String{Value: l + r}

	case "==":
		return btoo(l == r)

	case "!=":
		return btoo(l != r)

	default:
		lt := left.Type().String()
		rt := right.Type().String()
		return newError("invalid operator: %s %s %s", lt, op, rt)
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

func evalIfExpr(ie *ast.IfExpression, env *obj.Env) obj.Object {
	var cond = Eval(ie.Condition, env)

	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func isTruthy(cond obj.Object) bool {
	return cond != NULL && cond != FALSE
}

func evalBlockStatement(block *ast.BlockStatement, env *obj.Env) obj.Object {
	var res obj.Object

	for _, s := range block.Statements {
		res = Eval(s, env)

		if res != nil {
			rt := res.Type()
			if rt == obj.RETURN || rt == obj.ERROR {
				return res
			}
		}
	}
	return res
}

func evalIdentifier(node *ast.Identifier, env *obj.Env) obj.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if val, ok := builtins[node.Value]; ok {
		return val
	}

	return newError("identifier not found: %s", node.Value)
}

func evalIndexExpression(left, index obj.Object) obj.Object {
	switch {
	case left.Type() == obj.ARRAY && index.Type() == obj.INT:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index obj.Object) obj.Object {
	var arr = array.(*obj.Array)
	var idx = index.(*obj.Integer).Value
	var max = int64(len(arr.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arr.Elements[idx]
}

func evalExpressions(exps []ast.Expression, env *obj.Env) []obj.Object {
	var ret []obj.Object

	for _, e := range exps {
		val := Eval(e, env)
		if isError(val) {
			return []obj.Object{val}
		}
		ret = append(ret, val)
	}
	return ret
}

func applyFunction(fn obj.Object, args []obj.Object) obj.Object {
	switch fn := fn.(type) {

	case *obj.Function:
		extEnv := extendFuncEnv(fn, args)
		result := Eval(fn.Body, extEnv)
		return unwrapReturnValue(result)

	case *obj.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type().String())
	}
}

func unwrapReturnValue(o obj.Object) obj.Object {
	if rv, ok := o.(*obj.ReturnValue); ok {
		return rv.Value
	}
	return o
}

func extendFuncEnv(fn *obj.Function, args []obj.Object) *obj.Env {
	var env = obj.NewEnclosedEnv(fn.Env)

	for i, p := range fn.Params {
		env.Set(p.Value, args[i])
	}
	return env
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

func Eval(node ast.Node, env *obj.Env) obj.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expr, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}

	case *ast.Boolean:
		return btoo(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpr(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &obj.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Params
		body := node.Body
		return &obj.Function{Params: params, Env: env, Body: body}

	case *ast.CallExpression:
		fn := Eval(node.Func, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &obj.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}
