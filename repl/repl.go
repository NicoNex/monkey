package repl

import (
	"fmt"
	"github.com/NicoNex/monkey/evaluator"
	"github.com/NicoNex/monkey/lexer"
	"github.com/NicoNex/monkey/obj"
	"github.com/NicoNex/monkey/parser"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func printParserErrors(errs []string, out io.Writer) {
	for _, e := range errs {
		fmt.Fprintln(out, e)
	}
}

func Run() {
	var env = obj.NewEnv()
	var initState *terminal.State

	initState, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer terminal.Restore(0, initState)

	term := terminal.NewTerminal(os.Stdin, ">>> ")
	for {
		input, err := term.ReadLine()
		if err != nil {
			// Quit without error on Ctrl^D.
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}

		tokens := lexer.Lex(input)
		p := parser.New(tokens)
		prog := p.Parse()

		if errs := p.Errors(); len(errs) != 0 {
			printParserErrors(errs, term)
			continue
		}

		if val := evaluator.Eval(prog, env); val != nil {
			fmt.Fprintln(term, val.Inspect())
		}
	}
}
