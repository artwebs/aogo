package web

type Namespace struct {
	prefix   string
	handlers *ControllerRegistor
}
type innnerNamespace func(*Namespace)

func NewNamespace(prefix string, params ...innnerNamespace) *Namespace {
	ns := &Namespace{
		prefix:   prefix,
		handlers: NewControllerRegistor(),
	}
	for _, p := range params {
		p(ns)
	}
	return ns
}

// Namespace add sub Namespace
func NSNamespace(prefix string, params ...innnerNamespace) innnerNamespace {
	return func(ns *Namespace) {
		n := NewNamespace(prefix, params...)
		ns.Namespace(n)
	}
}

// register Namespace into beego.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		for k, v := range n.handlers.routers {
			if t, ok := BeeApp.Handlers.routers[k]; ok {
				addPrefix(v, n.prefix)
				BeeApp.Handlers.routers[k].AddTree(n.prefix, v)
			} else {
				t = NewTree()
				t.AddTree(n.prefix, v)
				addPrefix(t, n.prefix)
				BeeApp.Handlers.routers[k] = t
			}
		}
		if n.handlers.enableFilter {
			for pos, filterList := range n.handlers.filters {
				for _, mr := range filterList {
					t := NewTree()
					t.AddTree(n.prefix, mr.tree)
					mr.tree = t
					BeeApp.Handlers.insertFilterRouter(pos, mr)
				}
			}
		}
	}
}
