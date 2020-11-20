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
	FUNCTION
)

var typrepr = map[Type]string{
	NULL:     "NULL",
	ERROR:    "ERROR",
	INT:      "INTEGER",
	BOOL:     "BOOLEAN",
	RETURN:   "RETURN",
	FUNCTION: "FUNCTION",
}

func (t Type) String() string {
	return typrepr[t]
}
