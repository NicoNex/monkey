package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

func Start(in io.Reader, out io.Writer) {
	var scanner = bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, ">>> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		tokens := lexer.Lex(line)

		for t := range tokens {
			if t.Typ != token.EOF {
				fmt.Fprintln(out, t)
			}
		}
	}
}
