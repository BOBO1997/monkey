package ast

import (
	"bytes"

	"github.com/BOBO1997/monkey/token"
)

// Node interface requires TokenLiteral method
/* structs supporting Node interface:
Program
Identifier
LetStatement
*/
type Node interface {
	TokenLiteral() string // used only when debug and test
	String() string
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
	tokenLiteral := ""
	if len(p.Statements) > 0 {
		tokenLiteral = p.Statements[0].TokenLiteral()
	}
	return tokenLiteral
}

// String method of Program struct
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement is a struct for "let" statement
// "let" is a statement with identifier and expression
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // note: identifier is a struct
	Value Expression  // note: Expression is an interface
}

// StatementNode method of LetStatement struct,
func (ls *LetStatement) StatementNode() {}

// TokenLiteral method of LetStatement struct
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String method of LetStatement struct
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement is a struct
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression  // note: Expression is an interface
}

// StatementNode method of ReturnStatement struct,
func (rs *ReturnStatement) StatementNode() {}

// TokenLiteral method of ReturnStatement struct
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String method of ReturnStatement struct
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is a struct
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression // note : Expression is an interface
}

// StatementNode method of ExpressionStatement struct,
func (es *ExpressionStatement) StatementNode() {}

// TokenLiteral method of ExpressionStatement struct
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String method of Expressiontatement struct
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Identifier is a structure for token.Ident
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

// expressionNode method of Indentifier struct
func (i *Identifier) expressionNode() {}

// TokenLiteral method of Identifier struct
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
