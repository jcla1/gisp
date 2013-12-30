package main

import "fmt"

var Builtins = map[Symbol]ApplyFn{
	"lambda": lambda,
	"macro":  newMacro,
	"define": define,
	"assert": assert,
	"let":    let,
	"set":    set,
	"begin":  begin,
	"apply":  apply,
	"quote":  quote,
	// "true?": True,
	"println": println,
	"if":      iff,
	"dec":     dec,
	"inc":     inc,
	"+":       plus,
	"*":       times,
	"==":      eq,
}

func lambda(s *Scope, args []Any) Any {
	c := &closure{s, args[0].([]Any), args[1:]}
	return c
}

func newMacro(s *Scope, args []Any) Any {
	return NewMacro(&closure{NewScope(s), args[0].([]Any), args[1:]})
}

func define(s *Scope, args []Any) Any {
	err := s.env.Add(args[0].(Symbol), s.Eval(args[1]))
	if err != nil {
		panic(err)
	}
	return nil
}

func assert(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	for _, v := range args {
		if True(s, []Any{v}) != true {
			panic("assert failed")
		}
	}
	return nil
}

func set(s *Scope, args []Any) Any {
	s.env[args[0].(Symbol)] = s.Eval(args[1])
	return nil
}

func let(s *Scope, args []Any) Any {
	for _, v := range args[0].([]Any) {
		set(s, v.([]Any))
	}
	return s.Eval(args[len(args)-1])
}

func begin(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	return args[len(args)-1]
}

func inc(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	return args[0].(int64) + 1
}

func dec(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	return args[0].(int64) - 1
}

func apply(s *Scope, args []Any) Any {
	f := s.Eval(args[0])
	arguments := s.Eval(args[1]).([]Any)

	switch fn := f.(type) {
	case Function:
		return fn.Apply(arguments)
	case ApplyFn:
		return fn(s, arguments)
	default:
		panic(fmt.Errorf("Not a function (apply): %s", args[0]))
	}
}

func quote(s *Scope, args []Any) Any {
	return args[0]
}

func True(s *Scope, args []Any) Any {
	switch x := args[0].(type) {
	case []Any:
		if len(x) < 1 {
			return false
		} else {
			for _, v := range x {
				if v.(bool) == false {
					return false
				}
			}
		}
	case bool:
		return x
	}

	return true
}

func println(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	for _, v := range args {
		fmt.Println("->>", v)
	}
	return nil
}

func iff(s *Scope, args []Any) Any {
	if True(s, s.EvalAll(args[:1])) == true {
		return s.Eval(args[1])
	} else {
		return s.Eval(args[2])
	}
}

func plus(s *Scope, args []Any) Any {
	var c int64 = 0
	args = s.EvalAll(args)
	for _, n := range args {
		c += n.(int64)
	}
	return c
}

func times(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	var c int64 = 1
	for _, n := range args {
		c *= n.(int64)
	}

	return c
}

func eq(s *Scope, args []Any) Any {
	args = s.EvalAll(args)

	for i := 1; i < len(args); i++ {
		if args[i] != args[i-1] {
			return false
		}
	}
	return true
}
