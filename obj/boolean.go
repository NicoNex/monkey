package obj

import "strconv"

// TODO: consider refactoring this into 'type Boolean bool' if doable.
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

func (b *Boolean) Type() Type {
	return BOOL
}
