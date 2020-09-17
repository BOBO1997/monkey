package ast

import (
	"bytes"
	"strings"

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
		return es.Expression.String() // + ";"
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

// String method of Identifier struct
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral is a struct for token.Int
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// expressionNode method of IntegerLiteral struct
func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral method of IntegerLiteral struct
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String method of IntegerLiteral struct
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// StringLiteral is a struct for token.String
type StringLiteral struct {
	Token token.Token
	Value string
}

// expressionNode method of StringLiteral struct
func (sl *StringLiteral) expressionNode() {}

// TokenLiteral method of StringLiteral struct
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String method of StringLiteral struct
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// ArrayLiteral is a struct for token.Array
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

// expressionNode method of ArrayLiteral struct
func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral method of ArrayLiteral struct
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String method of ArrayLiteral struct
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, element := range al.Elements {
		elements = append(elements, element.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// expressions

// index expression

// IndexExpression struct
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

// expressionNode method of IndexExpression struct
func (ie *IndexExpression) expressionNode() {}

// TokenLiteral method of IndexExpression struct
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String method of IndexExpression struct
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	return out.String()
}

// HashLiteral is a struct for token.Hash
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

// expressionNode method of HashLiteral struct
func (hl *HashLiteral) expressionNode() {}

// TokenLiteral method of HashLiteral struct
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

// String method of ArrayLiteral struct
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, item := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+item.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// prefix expression

// PrefixExpression is a struct
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// expressionNode method of PrefixExpression struct
func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral method of PrefixExpression struct
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String method of PrefixExpression struct
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression is a struct
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// expressionNode method of InfixExpression struct
func (ie *InfixExpression) expressionNode() {}

// TokenLiteral method of InfixExpression struct
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String method of InfixExpression struct
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean is a struct
type Boolean struct {
	Token token.Token
	Value bool
}

// expressionNode method of Boolean struct
func (b *Boolean) expressionNode() {}

// TokenLiteral method of Boolean struct
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String method of Boolean struct
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression is a struct
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// expressionNode method of IfExpression struct
func (ife *IfExpression) expressionNode() {}

// TokenLiteral method of IfExpression struct
func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

// String method of IfExpression struct
func (ife *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	//out.WriteString(" then ")
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())
	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.String())
	}
	return out.String()
}

// BlockStatement is a struct
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

// expressionNode method of BlockStatement struct
func (bs *BlockStatement) expressionNode() {}

// TokenLiteral method of BlockStatement struct
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String method of BlockStatement struct
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	//out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	//out.WriteString("}")
	return out.String()
}

// FunctionLiteral is a struct
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

// expressionNode method of FunctionLiteral struct
func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral method of FunctionLiteral struct
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String method of FunctionLiteral struct
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression is a struct
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

// expressionNode method of CallExpression struct
func (ce *CallExpression) expressionNode() {}

// TokenLiteral method of CallExpression struct
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String method of CallExpression struct
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}

// MacroLiteral is a struct
type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

// expressionNode method of MacroLiteral struct
func (ml *MacroLiteral) expressionNode() {}

// TokenLiteral method of MacroLiteral struct
func (ml *MacroLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

// String method of MacroLiteral struct
func (ml *MacroLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(ml.Body.String())
	return out.String()
}
