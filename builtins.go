package main

import (
	"fmt"
)

var Builtins = map[Symbol]ApplyFn{
	"quote":      quote,
	"quasiquote": quasiquote,
	"lambda":     lambda,
	"macro":      newMacro,
	"macrox":     macroExpand,
	"define":     define,
	"if":         iff,
	"set!":       set,
	"let":        let,
	// wrapped funcs
	"==":      equality,
	"+":       addition,
	"-":       subtraction,
	"*":       multiplication,
	"car":     car,
	"cdr":     cdr,
	"eval":    eval,
	"true?":   True,
	"assert":  assert,
	"begin":   begin,
	"println": println,
}

// The "u" stands for unwrapped
var (
	eval           = wrap(wrap(uEval))
	equality       = wrap(uEquality)
	addition       = wrap(uAddition)
	subtraction    = wrap(uSubtraction)
	multiplication = wrap(uMultiplication)
	car            = wrap(uCar)
	cdr            = wrap(uCdr)
	True           = wrap(uTrue)
	assert         = wrap(uAssert)
	begin          = wrap(uBegin)
	println        = wrap(uPrintln)
)

func wrap(fn ApplyFn) ApplyFn {
	return func(s *Scope, args []Any) Any {
		args = s.EvalAll(args)
		return fn(s, args)
	}
}

func quote(s *Scope, args []Any) Any {
	return args[0]
}

func quasiquote(s *Scope, args []Any) Any {

	switch arg := args[0].(type) {
	case []Any:
		return resolveUnquoteSplices(s, resolveUnquotes(s, arg).([]Any))
	}
	return args[0]
}

func resolveUnquotes(s *Scope, sexp []Any) Any {
	if len(sexp) < 1 {
		return sexp
	}

	if sexp[0] == Symbol("unquote") {
		return s.Eval(sexp[1])
	} else if sexp[0] == Symbol("quasiquote") {
		return sexp
	}

	newSexp := make([]Any, len(sexp))
	for i, val := range sexp {
		switch val := val.(type) {
		case []Any:
			newSexp[i] = resolveUnquotes(s, val)
		default:
			newSexp[i] = val
		}
	}

	return newSexp
}

func resolveUnquoteSplices(s *Scope, sexp []Any) Any {
	if len(sexp) < 1 {
		return sexp
	}

	if sexp[0] == Symbol("unquote-splice") {
		return s.Eval(sexp[1])
	} else if sexp[0] == Symbol("quasiquote") {
		return sexp
	}

	for i := 0; i < len(sexp); i++ {
		val := sexp[i]
		switch val := val.(type) {
		case []Any:
			if len(val) > 1 && val[0] == Symbol("unquote-splice") {
				sexp = append(sexp[:i], append(resolveUnquoteSplices(s, val).([]Any), sexp[i+1:]...)...)
			} else {
				sexp[i] = resolveUnquoteSplices(s, val)
			}

		}
	}

	return sexp
}

func lambda(s *Scope, args []Any) Any {
	return &closure{s, args[0].([]Any), args[1:]}
}

func newMacro(s *Scope, args []Any) Any {
	return &macro{nil, args[0].([]Any), args[1:]}
}

func macroExpand(s *Scope, args []Any) Any {
	args = args[0].([]Any)
	return internalMacroExpand(s, args)
}

func internalMacroExpand(s *Scope, args []Any) Any {
	m := args[0]

	// The macro is guaranteed to be defined
	// since we s.Eval its Symbol beforehand
	m, _ = s.Lookup(m.(Symbol))

	return m.(Macro).Apply(s, args[1:])
}

func define(s *Scope, args []Any) Any {
	err := s.Add(args[0].(Symbol), s.Eval(args[1]))
	if err != nil {
		panic(err)
	}

	return nil
}

func iff(s *Scope, args []Any) Any {
	if True(s, []Any{args[0]}) == true {
		return s.Eval(args[1])
	} else {
		if len(args) > 2 {
			return s.Eval(args[2])
		} else {
			return nil
		}
	}
}

func set(s *Scope, args []Any) Any {
	s.Override(args[0].(Symbol), s.Eval(args[1]))
	return nil
}

func let(s *Scope, args []Any) Any {
	for _, v := range args[0].([]Any) {
		set(s, v.([]Any))
	}
	res := s.EvalAll(args[1:])
	return res[len(res)-1]
}

// Wrapped functions
// -----------------

func uEval(s *Scope, args []Any) Any {
	return args[len(args)-1]
}

func uEquality(s *Scope, args []Any) Any {
	for i := 1; i < len(args); i++ {
		if args[i-1] != args[i] {
			return false
		}
	}
	return true
}

func uAddition(s *Scope, args []Any) Any {
	var c int64 = 0
	for _, val := range args {
		c += val.(int64)
	}
	return c
}

func uSubtraction(s *Scope, args []Any) Any {
	var c int64 = args[0].(int64)
	for i := 1; i < len(args); i++ {
		c -= args[i].(int64)
	}
	return c
}

func uMultiplication(s *Scope, args []Any) Any {
	var c int64 = args[0].(int64)
	for i := 1; i < len(args); i++ {
		c *= args[i].(int64)
	}
	return c
}

func uCar(s *Scope, args []Any) Any {
	return args[0]
}

func uCdr(s *Scope, args []Any) Any {
	return args[1:]
}

func uTrue(s *Scope, args []Any) Any {
	for _, val := range args {
		if val == nil || val == false {
			return false
		}
	}
	return true
}

func uAssert(s *Scope, args []Any) Any {
	for _, val := range args {
		if True(s, []Any{val}) != true {
			panic("assertion failed")
		}
	}
	return nil
}

func uBegin(s *Scope, args []Any) Any {
	return args[len(args)-1]
}

func uPrintln(s *Scope, args []Any) Any {
	for _, val := range args {
		fmt.Println(val)
	}
	return nil
}
