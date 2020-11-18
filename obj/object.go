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
	RETURN_VALUE
)

var typrepr = map[Type]string{
	NULL:         "NULL",
	INT:          "INTEGER",
	BOOL:         "BOOLEAN",
	RETURN_VALUE: "RETURN_VALUE",
}

func (t Type) String() string {
	return typrepr[t]
}
