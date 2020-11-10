package main

import (
	"fmt"

	"monkey/lexer"
)

func main() {
	for t := range lexer.Lex(`=+(){},;`) {
		fmt.Println(t)
	}
}
