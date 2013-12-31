package main

import (
	"fmt"
)

var Builtins = map[Symbol]ApplyFn{
	"quote":  quote,
	"lambda": lambda,
	"define": define,
	"if":     iff,
	// wrapped funcs
	"==":      equality,
	"+":       addition,
	"-":       subtraction,
	"true?":   True,
	"assert":  assert,
	"println": println,
}

// The "u" stands for unwrapped
var (
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
