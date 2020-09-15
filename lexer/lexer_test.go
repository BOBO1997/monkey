package lexer

import (
	"testing"

	"github.com/BOBO1997/monkey/token"
)

func TestNextToken1(t *testing.T) {
	input := `=+(){},;`

	tests := []struct { // Token struct for testing
		expectedType    token.TokenType // Type
		expectedLiteral string          // Literal
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	lex := New(input)          // make a new lexer
	for i, tt := range tests { // read token one by one
		tok := lex.NextToken() // update the counters

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected %q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokentype wrong. expected %q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;
		let add = fn (x, y) { x + y; };
		let result = add(five, ten);
		if five + ten < 20 {
			return 15;
		} else {
			return 20;
		}
		ten == five;
		10 != 5;
		"foobar"
		"foo bar"
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 0
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// 5
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		// 10
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		// 15
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		// 20
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		// 25
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		// 30
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		// 35
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.IDENT, "five"},
		{token.PLUS, "+"},
		{token.IDENT, "ten"},
		// 40
		{token.LT, "<"},
		{token.INT, "20"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "15"},
		// 45
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		// 50
		{token.INT, "20"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "ten"},
		{token.EQ, "=="},
		// 55
		{token.IDENT, "five"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "5"},
		//60
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.EOF, ""},
	}

	lex := New(input)
	for i, tt := range tests {
		tok := lex.NextToken()
		//fmt.Printf("%#U\n", lex.ch)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected %q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokenliteral wrong. expected %q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
