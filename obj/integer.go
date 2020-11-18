package obj

import "strconv"

// TODO: consider refactoring this into 'type Integer int64' if doable.
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Type() Type {
	return INT
}
