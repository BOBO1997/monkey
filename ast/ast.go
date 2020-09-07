package ast

import "github.com/BOBO1997/monkey/token"

// Node interface requires TokenLiteral method
/* structs supporting Node interface:
Program
Identifier
LetStatement
*/
type Node interface {
	TokenLiteral() string // used only when debug and test
}

// Statement interface requires Node interface and StatementNode method
type Statement interface {
	Node
	StatementNode()
}

// Expression interface requires Node interface and expressionNode method
type Expression interface {
	Node
	expressionNode()
}

// Program is a structof whole ast, which is relaized by a slice of Statement interface
type Program struct {
	Statements []Statement
}

// TokenLiteral method of Program struct, returns the token literal of the first statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement is a struct for "let" statement
// "let" is a statement with identifier and expression
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // pointer to child (identifier)
	Value *Expression // pointer to child (expression)
}

// StatementNode method of LetStatement struct,
func (ls *LetStatement) StatementNode() {}

// TokenLiteral method of LetStatement struct
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Identifier is a structure for token.Ident
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

// expressionNode method of Indentifier struct
func (i *Identifier) expressionNode() {}

// TokenLiteral method of Identifier struct
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
