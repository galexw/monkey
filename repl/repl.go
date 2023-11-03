package repl

import (
	"bufio"
	"fmt"
	"galexw/monkey/lexer"
	"galexw/monkey/token"
	"io"
)

const PROMPT = `>>> `

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)

		for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
