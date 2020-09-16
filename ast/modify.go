package ast

// ModifierFunc type
type ModifierFunc func(Node) Node

// Modify function
func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
	case *PrefixExpression:
	case *IndexExpression:
	case *IfExpression:
	case *BlockStatement:
	case *ReturnStatement:
	case *LetStatement:
	case *FunctionLiteral:
	case *ArrayLiteral:
	case *HashLiteral:
	}
	return modifier(node)
}
