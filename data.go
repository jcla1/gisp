package main

import (
	"fmt"
)

type Any interface{}
type Symbol string
type Environment map[Symbol]Any

type Scope struct {
	env    Environment
	parent *Scope
}

type macro struct {
	t string
	c *closure
}

func NewMacro(c *closure) macro {
	return macro{"", c}
}

func (m macro) String() string {
	return "#<macro>"
}

func (m macro) Apply(args []Any) Any {
	for i, v := range m.c.vars {
		m.c.s.env[v.(Symbol)] = args[i]
	}

	evaled := m.c.s.EvalAll(m.c.body)
	return evaled[len(evaled)-1]
}

func NewRootScope() *Scope {
	s := &Scope{make(Environment), nil}
	for k, v := range Builtins {
		s.env[k] = v
	}
	return s
}

func NewScope(parent *Scope) *Scope {
	return &Scope{make(Environment), parent}
}

func (s *Scope) macroExpand(sexp Any) Any {
	switch sexp := sexp.(type) {
	case []Any:
		f := sexp[0]
		switch f := f.(type) {
		case Symbol:
			m, err := s.Lookup(f)
			if err != nil {
				panic(err)
			}

			_, ok := m.(macro)
			if !ok {
				return sexp
			} else {
				fmt.Println(m.(macro).Apply(sexp[1:]))
			}
			return sexp
		default:
			return sexp
		}
	default:
		return sexp
	}

	panic("can't get here")
}

func (s *Scope) Eval(sexp Any) Any {
	// sexp = s.macroExpand(sexp)
	switch sexp := sexp.(type) {
	case []Any:
		if len(sexp) < 1 {
			return nil
		} else {
			return s.evalFunctionCall(sexp)
		}
	case Symbol:
		v, err := s.Lookup(sexp)
		if err != nil {
			panic(err)
		}
		return v
	default:
		return sexp
	}

	return nil
}

func (s *Scope) EvalAll(sexps []Any) []Any {
	res := make([]Any, len(sexps))
	for i, v := range sexps {
		res[i] = s.Eval(v)
	}

	return res
}

func (s *Scope) evalFunctionCall(sexp []Any) Any {
	f := s.Eval(sexp[0])

	switch fn := f.(type) {
	case Function:
		return fn.Apply(sexp[1:])
	case ApplyFn:
		return fn(s, sexp[1:])
	default:
		panic(fmt.Errorf("Not a function (main): %s", sexp[0]))
	}
}

func (e Environment) Add(s Symbol, val Any) error {
	_, ok := e[Symbol(s)]

	if ok {
		return fmt.Errorf("symbol: \"%s\" already defined in this scope", s)
	}

	e[s] = val
	return nil
}

func (s *Scope) Lookup(sym Symbol) (Any, error) {
	val, ok := s.env[sym]

	if !ok && s.parent != nil {
		return s.parent.Lookup(sym)
	} else if !ok && s.parent == nil {
		return nil, fmt.Errorf("unknown variable: %s", sym)
	}

	return val, nil
}

type ApplyFn func(s *Scope, args []Any) Any
type Function interface {
	Apply([]Any) Any
}

type closure struct {
	s          *Scope
	vars, body []Any
}

func (c *closure) Apply(args []Any) Any {
	c.s = NewScope(c.s)
	args = c.s.EvalAll(args)
	for i, v := range c.vars {
		c.s.env[v.(Symbol)] = args[i]
	}

	evaled := c.s.EvalAll(c.body)
	return evaled[len(evaled)-1]
}

func (c *closure) String() string {
	return "#<closure>"
}
