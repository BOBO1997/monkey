package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/BOBO1997/monkey/lexer"
	"github.com/BOBO1997/monkey/token"
)

// PROMPT is a const string
const PROMPT = ">> "

// Start function scans one line and output the lexical tokens
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Printf(PROMPT)
	scanned := scanner.Scan()
	if !scanned {
		return
	}
	line := scanner.Text()
	lex := lexer.New(line)

	for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}
