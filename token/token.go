package token

// TokenType is a string representing the type of each token
type TokenType string

// Token is a struct which fields are
// Type : type of the token
// Literal : literal of the token, not the raw word
type Token struct {
	Type    TokenType
	Literal string
}

// token.TokenType is implemented by const string
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	LT  = "<"
	GT  = ">"
	LEQ = "<="
	GEQ = ">="
	EQ  = "=="
	NEQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	THEN     = "THEN"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
)

// token.keywords is a map which contains the reserved identifier
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"return": RETURN,
	"for":    FOR,
}

// LookupIdent function identify whether the identifier is keyword or not, and return its token type
func LookupIdent(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENT
}
