package evaluator

import (
	"fmt"

	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/object"
	"github.com/BOBO1997/monkey/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) { // do nothing
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		if len(call.Arguments) != 1 {
			return node
		}
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToAstNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		tok := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{
			Token: tok,
			Value: obj.Value,
		}
	case *object.Boolean:
		var tok token.Token
		if obj.Value {
			tok = token.Token{
				Type:    token.TRUE,
				Literal: "true",
			}
		} else {
			tok = token.Token{
				Type:    token.FALSE,
				Literal: "false",
			}
		}
		return &ast.Boolean{
			Token: tok,
			Value: obj.Value,
		}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
