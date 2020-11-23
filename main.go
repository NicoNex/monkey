package main

import (
	"fmt"
	"github.com/NicoNex/monkey/repl"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	repl.Run()
}
