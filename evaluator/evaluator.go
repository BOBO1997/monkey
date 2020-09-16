package evaluator

import (
	"fmt"

	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// BaseEnv is the base environment of repl
var BaseEnv *object.Environment

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: arg.Value[0:1]}
				}
				return NULL
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				length := len(arg.Value)
				if length > 0 {
					return &object.String{Value: arg.Value[length-1 : length]}
				}
				return NULL
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					return arg.Elements[length-1]
				}
				return NULL
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				length := len(arg.Value)
				if length > 0 {
					return &object.String{Value: arg.Value[1:length]}
				}
				return NULL
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					rest := make([]object.Object, length-1, length-1)
					copy(rest, arg.Elements[1:length])
					return &object.Array{Elements: rest}
				}
				return NULL
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want=2", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newArr := make([]object.Object, length+1, length+1)
				copy(newArr, arg.Elements)
				newArr[length] = args[1]
				return &object.Array{Elements: newArr}
			default:
				return newError("argument to `push` not supported, got %s", args[0].Type())
			}
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"__inspect_env__": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments, got=%d, want= 0 or 1", len(args))
			}
			if len(args) == 0 {
				fmt.Println(object.InspectEnvironment(BaseEnv))
			} else {
				switch arg := args[0].(type) {
				case *object.Function:
					fmt.Println(object.InspectEnvironment(arg.Env))
				}
			}
			return NULL
		},
	},
}

// Eval function
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters // raw ast
		body := node.Body         // raw ast
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression: // this can handle recursive function because the defined function is already in base env and the Parameters field and Body field of object.Function are hold by ast
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) { // if error occur
			return elements[0]
		}
		return &object.Array{Elements: elements}
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
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

// evalProgram function
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

// evalBlockStatement function
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// evalStatements function
func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if returnValue, ok := result.(*object.ReturnValue); ok { // ok; type assersion should have error checking
			return returnValue.Value
		}
	}
	return result
}

// nativeBoolToBooleanObject function
func nativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

// evalPrefixExpression function
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s %s", operator, right.Type())
	}
}

// evalBangOperatorExpression function
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusOperatorExpression function
// this function is only called when right is *object.Integer type
func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ { // ok; type assersion should have error checking
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value // ok; type assersion should have error checking
	return &object.Integer{Value: -value}
}

// evalInfixExpression function
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ: // evaluated first
		return evalIntegerInfixExpression(operator, left, right)
	/*
		case left.Type() != right.Type():
			return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	*/
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ: // evaluated first
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression function
// this function is only called when left and right are both *object.String type
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value   // ok; type assersion should have error checking
	rightVal := right.(*object.String).Value // ok; type assersion should have error checking
	return &object.String{Value: leftVal + rightVal}
}

// evalIntegerInfixExpression function
// this function is only called when left and right are both *object.Integer type
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value   // ok; type assersion should have error checking
	rightVal := right.(*object.Integer).Value // ok; type assersion should have error checking
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIfExpression function
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValues(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return newError("not a function: %s", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValues(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok { // ok; type assersion should have error checking
		return returnValue.Value
	}
	return obj
}

// evalIndexExpression function evaluates index access expression
// index expression belongs to
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalApplyIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalApplyIndexExpression function accesses the index-th element of the array
func evalApplyIndexExpression(array, index object.Object) object.Object {
	arrayObject, ok := array.(*object.Array) // ok; type assersion should have error checking
	if !ok {
		return newError("not Array: got=%s", array.Type())
	}
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || max < idx {
		return NULL
	}
	return arrayObject.Elements[idx]
}

// evalHashLiteral function makes a hash object, which contains a map from object to object
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyExpr, itemExpr := range node.Pairs {
		key := Eval(keyExpr, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable) // ok; type assersion should have error checking
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		item := Eval(itemExpr, env)
		if isError(item) {
			return item
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: item}
	}
	return &object.Hash{Pairs: pairs}
}

// evalHashIndexExpression function
func evalHashIndexExpression(left, index object.Object) object.Object {
	hashObject, ok := left.(*object.Hash) // ok; type assersion should have error checking
	if !ok {
		return newError("not Hash: got=%s", left.Type())
	}
	key, ok := index.(object.Hashable) // ok; type assersion should have error checking
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
