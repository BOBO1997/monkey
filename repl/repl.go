package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/BOBO1997/monkey/evaluator"

	"github.com/BOBO1997/monkey/lexer"
	"github.com/BOBO1997/monkey/parser"
)

// PROMPT is a const string
const PROMPT = ">> "

// Start function scans one line and output the lexical tokens
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			continue
		}
		line := scanner.Text()
		if line == ":q" {
			fmt.Printf("Bye bye! \n")
			return
		}
		lex := lexer.New(line)
		p := parser.New(lex)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		// lexer output
		/*
			for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		*/
		// parser output
		/*
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		*/
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
