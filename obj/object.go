package obj

type Object interface {
	Type() Type
	Inspect() string
}

type Type int

const (
	NULL Type = iota
	INT
	BOOL
)

var typrepr = map[Type]string{
	NULL: "NULL",
	INT: "INTEGER",
	BOOL: "BOOLEAN",
}

func (t Type) String() string {
	return typrepr[t]
}
