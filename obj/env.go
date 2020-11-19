package obj

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	return &Env{store: make(map[string]Object)}
}

func (e *Env) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
