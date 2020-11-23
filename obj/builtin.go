package obj

type BuiltinFn func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() Type {
	return BUILTIN
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}
