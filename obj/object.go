package obj

type Object interface {
	Type() Type
	Inspect() string
}

type Type int

const (
	NULL Type = iota
	ERROR
	INT
	BOOL
	RETURN
)

var typrepr = map[Type]string{
	NULL:   "NULL",
	ERROR:  "ERROR",
	INT:    "INTEGER",
	BOOL:   "BOOLEAN",
	RETURN: "RETURN",
}

func (t Type) String() string {
	return typrepr[t]
}
