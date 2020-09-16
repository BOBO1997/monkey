package lexer

import (
	"github.com/BOBO1997/monkey/token"
)

// Lexer is a struct holding the information of whole source code and the counter of lexer
type Lexer struct {
	input        string // the whole source
	position     int    // the positing currently reading (alrerady read)
	readPosition int    // the next position to be read
	ch           rune   // one charactor at the position
}

// New function makes a new *Lexer struct
// input is a string of raw token
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar method update the *Lexer struct, with making value of ch field and updating position and readPosition
// this method is the updating function for lexer counters
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 represents EOF
	} else {
		l.ch = []rune(l.input)[l.readPosition] // read a new rune
	}
	l.position = l.readPosition
	l.readPosition++
}

// readIdentifier method reads forward the source code and return an identifier
func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// readNumber method reads forward the source code and return a number
func (l *Lexer) readNumber() string {
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// readString method reads forward the source code and return a string
func (l *Lexer) readString() string {
	startPos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[startPos:l.position]
}

// peekChar method peeks the next rune, for finding the operator with two rune
func (l *Lexer) peekChar() rune {
	if len([]rune(l.input)) <= l.readPosition {
		return 0
	} else {
		return rune(l.input[l.readPosition])
	}
}

// skipWhitespace method skips the white space and escape sequences
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken method returns a token.Token, according to the symbol in Lexer.ch
// This method is to update the counters of Lexer structure and return the token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = "EOF"
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// newToken function creates a new token.Token, with its token type and symbol.
func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)} // make and return a new token
}

// isLetter function checks whether the current rune is a letter or not
func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

// isDigit function checks whether the current rune is a digit or not
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
