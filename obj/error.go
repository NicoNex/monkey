package obj

import "fmt"

type Error struct {
	Msg string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("error: %s", e.Msg)
}
