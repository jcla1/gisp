package main

import (
	"fmt"
)

var Builtins = map[Symbol]ApplyFn{
	"quote":      quote,
	"quasiquote": quasiquote,
	"lambda":     lambda,
	"define":     define,
	"if":         iff,
	// wrapped funcs
	"==":      equality,
	"+":       addition,
	"-":       subtraction,
	"eval":    eval,
	"true?":   True,
	"assert":  assert,
	"println": println,
}

// The "u" stands for unwrapped
var (
	eval        = wrap(wrap(uEval))
	equality    = wrap(uEquality)
	addition    = wrap(uAddition)
	subtraction = wrap(uSubtraction)
	True        = wrap(uTrue)
	assert      = wrap(uAssert)
	println     = wrap(uPrintln)
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
		args[0] = resolveUnquoteSplices(s, resolveUnquotes(s, arg).([]Any))
	}
	return args[0]
}

func resolveUnquotes(s *Scope, sexp []Any) Any {
	if sexp[0] == Symbol("unquote") {
		return s.Eval(sexp[1])
	}

	for i, val := range sexp {
		switch val := val.(type) {
		case []Any:
			if len(val) > 0 && sexp[0] != Symbol("quasiquote") {
				sexp[i] = resolveUnquotes(s, val)
			}
		}
	}

	return sexp
}

func resolveUnquoteSplices(s *Scope, sexp []Any) Any {
	if sexp[0] == Symbol("unquote-splice") {
		return s.Eval(sexp[1])
	}

	for i, val := range sexp {
		switch val := val.(type) {
		case []Any:
			if len(val) > 0 && sexp[0] != Symbol("quasiquote") {
				a := resolveUnquoteSplices(s, val).([]Any)
				sexp = append(sexp[:i], append(a, sexp[i+1:]...)...)
			}
		}
	}

	return sexp
}

func lambda(s *Scope, args []Any) Any {
	// We don't assign a scope here, because
	// it is passed to fn.Apply when it is called
	return &closure{nil, args[0].([]Any), args[1:]}
}

func define(s *Scope, args []Any) Any {
	s.Add(args[0].(Symbol), s.Eval(args[1]))
	return nil
}

func iff(s *Scope, args []Any) Any {
	if True(s, []Any{args[0]}) == true {
		return s.Eval(args[1])
	} else {
		return s.Eval(args[2])
	}
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

func uPrintln(s *Scope, args []Any) Any {
	for _, val := range args {
		fmt.Println(val)
	}
	return nil
}
