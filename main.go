package main

import (
	"fmt"

	"monkey/lexer"
)

func main() {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);`

	for t := range lexer.Lex(input) {
		fmt.Println(t)
	}
}
