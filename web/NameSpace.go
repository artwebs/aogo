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

// nest Namespace
// usage:
//ns := beego.NewNamespace(“/v1”).
//Namespace(
//    beego.NewNamespace("/shop").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("shopinfo"))
//    }),
//    beego.NewNamespace("/order").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("orderinfo"))
//    }),
//    beego.NewNamespace("/crm").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("crminfo"))
//    }),
//)

// same as beego.Rourer
// refer: https://godoc.org/github.com/astaxie/beego#Router
func (n *Namespace) Router(rootpath string, c ControllerInterface, method string) *Namespace {
	n.handlers.routes[n.prefix+rootpath] = &Handler{controller: c, method: method}
	return n
}

// Namespace Router
func NSRouter(rootpath string, c ControllerInterface, method string) innnerNamespace {
	return func(ns *Namespace) {
		ns.Router(rootpath, c, method)
	}
}

func (n *Namespace) Namespace(ns ...*Namespace) *Namespace {
	for _, ni := range ns {
		for k, v := range ni.handlers.routes {
			register.routes[n.prefix+k] = v
		}
	}
	return n
}

func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		for k, v := range n.handlers.routes {
			register.routes[k] = v
		}
	}
}

// register Namespace into beego.Handler
// support multi Namespace
// func AddNamespace(nl ...*Namespace) {
// 	for _, n := range nl {
// 		for k, v := range n.handlers.routers {
// 			if t, ok := BeeApp.Handlers.routers[k]; ok {
// 				addPrefix(v, n.prefix)
// 				BeeApp.Handlers.routers[k].AddTree(n.prefix, v)
// 			} else {
// 				t = NewTree()
// 				t.AddTree(n.prefix, v)
// 				addPrefix(t, n.prefix)
// 				BeeApp.Handlers.routers[k] = t
// 			}
// 		}
// 		if n.handlers.enableFilter {
// 			for pos, filterList := range n.handlers.filters {
// 				for _, mr := range filterList {
// 					t := NewTree()
// 					t.AddTree(n.prefix, mr.tree)
// 					mr.tree = t
// 					BeeApp.Handlers.insertFilterRouter(pos, mr)
// 				}
// 			}
// 		}
// 	}
// }
