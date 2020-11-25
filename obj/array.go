package obj

import (
	"fmt"
	"strings"
)

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type {
	return ARRAY
}

func (a *Array) Inspect() string {
	var elements []string

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
