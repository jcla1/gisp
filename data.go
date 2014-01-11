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

func NewRootScope() *Scope {
	s := NewScope(nil)
	for k, v := range Builtins {
		s.env[k] = v
	}
	return s
}

func NewScope(parent *Scope) *Scope {
	return &Scope{make(Environment), parent}
}

func (s *Scope) evalFunctionCall(sexp []Any) Any {
	f := s.Eval(sexp[0])

	switch fn := f.(type) {
	case Macro:
		expansion := internalMacroExpand(s, sexp)
		return s.Eval(expansion)
	case Function:
		return fn.Apply(s, sexp[1:])
	case ApplyFn:
		return fn(s, sexp[1:])
	default:
		panic(fmt.Errorf("Not a function (main): %s", sexp[0]))
	}
}

func (s *Scope) Eval(sexp Any) Any {

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

func (s *Scope) Add(sym Symbol, val Any) error {
	_, ok := s.env[sym]

	if ok {
		return fmt.Errorf("symbol: \"%s\" already defined in this scope", s)
	}

	s.env[sym] = val
	return nil
}

func (s *Scope) Override(sym Symbol, val Any) {
	s.env[sym] = val
}

func (s *Scope) getRootScope() *Scope {
	if s.parent == nil {
		return s
	} else {
		return s.parent.getRootScope()
	}
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
