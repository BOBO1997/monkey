package parser

import (
	"testing"

	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/lexer"
)

// TestLetStatements function tests a sequence of let statements
func TestLetSatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 124455;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram() // parse whole program
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement does note contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) { // if false
			return
		}
	}
}

// testLetStatement function tests one let statement
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement) // type assersion
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got %s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.TokenLiteral())
	}
	return true
}

// TestReturnStatements return
func TestReturnStatements(t *testing.T) {
	input := `
		return 10;
		return 1342;
		return 2224111;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram() // parse whole program
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement does note contain 3 statements. got=%d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

// checkParserErrors function outputs the whole errors accumulated in Parser
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error : %q", msg)
	}
	t.FailNow()
}
