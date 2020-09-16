package parser

import (
	"fmt"
	"strconv"

	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/lexer"
	"github.com/BOBO1997/monkey/token"
)

// Parser is a struct for parsing whole program
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression               // prefix parse function
	infixParseFn  func(ast.Expression) ast.Expression // infix parse function
)

// New function creates a parser from lexer of whole program
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LEQ, p.parseInfixExpression)
	p.registerInfix(token.GEQ, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.nextToken() // go forward
	p.nextToken() // go forward
	return p
}

// Errors method of Parser returns errors field
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken method of Parser struct contains current token and next token by peeking
// directly operates on curToken field and peekToken field
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// peekError method of Parser adds error message to errors field if token type is not correct
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// ParseProgram method of Parser struct parses whole program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement method of Parser parses one statement
func (p *Parser) parseStatement() ast.Statement { // note: ast.Statement is an interface
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement method of Parser struct parses a let statement
// let statement is expected to be "let <identifier> = <expression>"
func (p *Parser) parseLetStatement() *ast.LetStatement { // note: ast.LetStatement is a struct
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: expression
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseReturnStatement method of Parser struct parses a return statement
// return statement is expected to be "return <expression>"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: expression
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement method of Parser struct
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

/* ====== expression ====== */
// parser functions for expression

// priority depth of each operator
const (
	_           int = iota
	LOWEST          // 0
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // + or -
	PRODUCT         // * or /
	PREFIX          // -X or !X
	CALL            // myFunction(X)
	INDEX           // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LEQ:      LESSGREATER,
	token.GEQ:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// parseExpressionStatement method of Parser struct parses expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression method of Parser struct parses expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type] // prefix is a function, search "prefix" at first, prefix in this case means the first operand
	if prefix == nil {                          // error : no such prefix
		p.noParsingPrefixFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix() // parseSomething(), including identifier, integer literal, prefix expression, ...

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // if the right operand is stronger: then return combined infix expression
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp // finish parsing, since infix operator is not found
		}
		p.nextToken()
		leftExp = infix(leftExp) // leftExp + operator + rightExp
	}
	return leftExp // leftExp
}

// parseIdentifier method of Parser struct returns ast.Expression interface, which contains an identifier
// identifier expression is expected to be "<identifier>;"
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral method of Parser struct returns ast.Expression interface, which contains an integer literal
// integer expression is expected to be "<integer literal>;"
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	literal := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64) // change string to int
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	literal.Value = value // value is int
	return literal
}

// prefix

// registerPrefix method of Parser struct
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// parsePrefixExpression method of Prefix struct returns ast.Expression interface, which contains prefix operator
// prefix expression is expected to be "<prefix operator> <expression>;"
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExoression"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX) // recursive call of parseExpression
	return expression
}

// noParsingPrefixFnError method of Parser struct stores an error message for no prefix error
func (p *Parser) noParsingPrefixFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseBoolean method of Parser struct
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

// parseStringLiteral method of Parser struct
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseExpressionList method of Parser struct
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// parseArrayLiteral method of Parser struct
func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.curToken,
		Elements: p.parseExpressionList(token.RBRACKET),
	}
}

// infix

// registerInfix method of Parser struct
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// peekPrecedence method of Parser struct
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence method of Parser struct
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parseInfixExprepssion method of Parser struct
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression { // ex) "5 + 10", left is 5, curToken is +, and right is 5
	// defer untrace(trace("parseInfixExoression"))
	expression := &ast.InfixExpression{ // already known
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence) // recursive call
	return expression
}

// grouped exression

// parseGroupedExpression method of Parser srruct
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST) // recursive call
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// if expression
// if expression is expected to be "if (<condition>) <consequence> else <alternative>;"
// <condition> is ast.Expression
// <consequence> is *ast.BlockStatement
// <alternative> is *ast.BlockStatement

// parseIfExpression method of Parser struct
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

// function literal
// function literal is expected to be "fn(<parameters>) <body>;"
// <parameters> is []*ast.Identifier
// <body> is *ast.BlockStatement

// parseFunctionLiteral method of Parser struct
func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	literal.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	literal.Body = p.parseBlockStatement()
	return literal
}

// parseFunctionParameters method of Parser struct
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() // skip ")"
		return params // empty
	}
	p.nextToken()

	param := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	params = append(params, param)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","
		param := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		params = append(params, param)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
}

// function call

// parseCallExpression method of Parser struct
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	// exp.Arguments = p.parseCallArguments()
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

/*
// parseCallArguments method of Parser struct
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken() // move to exp
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to ","
		p.nextToken() // skip ","
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
*/

// parseIndexExpression method of Parser struct
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

/* ====== assertion functions ====== */

// peekTokenIs method of Parse, checking the type of current token
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// curTokenIs method of Parser, checking the type of peeked token
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek method of Parser
// An assertion function of Parser, calling nextToken when the type of peeked token corresponds to what is expected to come next.
func (p *Parser) expectPeek(t token.TokenType) bool {
	expect := true
	if p.peekTokenIs(t) {
		p.nextToken() // go forward
	} else {
		p.peekError(t)
		expect = false
	}
	return expect
}
