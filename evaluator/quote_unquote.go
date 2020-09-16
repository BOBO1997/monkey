package evaluator

import (
	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
