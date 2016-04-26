package web

import (
	"github.com/artwebs/aogo/utils"
	"path"
	"reflect"
	"strings"
)

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

// Namespace Router
func NSAutoRouter(rootpath string, c ControllerInterface) innnerNamespace {
	return func(ns *Namespace) {

		reflectVal := reflect.ValueOf(c)
		rt := reflectVal.Type()
		ct := reflect.Indirect(reflectVal).Type()
		controllerName := strings.TrimSuffix(ct.Name(), "Controller")
		for i := 0; i < rt.NumMethod(); i++ {
			if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
				pattern := path.Join(rootpath, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
				ns.Router(pattern, c, rt.Method(i).Name)
			}

		}
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
			if !utils.InSlice(n.prefix, register.namespaces) {
				register.namespaces = append(register.namespaces, n.prefix)
			}
		}
	}
}
