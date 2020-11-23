package obj

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}
