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
	STRING
	RETURN
	FUNCTION
	BUILTIN
	ARRAY
)

var typrepr = map[Type]string{
	NULL:     "NULL",
	ERROR:    "ERROR",
	INT:      "INTEGER",
	BOOL:     "BOOLEAN",
	STRING:   "STRING",
	RETURN:   "RETURN",
	FUNCTION: "FUNCTION",
	BUILTIN:  "BUILTIN",
	ARRAY:    "ARRAY",
}

func (t Type) String() string {
	return typrepr[t]
}
