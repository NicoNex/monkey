package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

func printParserErrors(errs []string, out io.Writer) {
	for _, e := range errs {
		fmt.Fprintln(out, e)
	}
}

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
		p := parser.New(tokens)
		prog := p.Parse()

		if errs := p.Errors(); len(errs) != 0 {
			printParserErrors(errs, out)
			continue
		}
		fmt.Fprintln(out, prog.String())
	}
}
