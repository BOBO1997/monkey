package parser

import (
	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/lexer"
	"github.com/BOBO1997/monkey/token"
)

// Parser is a struct for parsing whole program
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

// New function
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

// Errors method of Parser returns p.errors
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken method of Parser struct contains current token and next token by peeking
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
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
	default:
		return nil
	}
}

// parseLetStatement method of Parser
func (p *Parser) parseLetStatement() *ast.LetStatement { // note: ast.LetStatement is a struct
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// peekTokenIs method of Parse, checking the type of current token
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// curTokenIs method of Parser, checking the type of peeked token
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek method of Parser
// An assertion function of Parser, calling nextToken when the type of peeked token corresponds to what is expected.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
