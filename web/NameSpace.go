package web

type NameSpace struct {
	prefix   string
	handlers *ControllerRegistor
}
type innnerNamespace func(*Namespace)

func NewNamespace(prefix string, params ...innnerNamespace) *Namespace {
	ns := &Namespace{
		prefix:   prefix,
		handlers: NewControllerRegister(),
	}
	for _, p := range params {
		p(ns)
	}
	return ns
}
