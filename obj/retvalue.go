package obj

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() Type {
	return RETURN_VALUE
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
