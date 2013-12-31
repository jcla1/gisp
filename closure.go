package main

type ApplyFn func(s *Scope, args []Any) Any
type Function interface {
	Apply(*Scope, []Any) Any
}

type closure struct {
	s          *Scope
	vars, body []Any
}

func (c *closure) Apply(s *Scope, args []Any) Any {
	args = s.EvalAll(args)
	c.s = NewScope(s)
	for i, v := range c.vars {
		c.s.env[v.(Symbol)] = args[i]
	}

	evaled := c.s.EvalAll(c.body)
	return evaled[len(evaled)-1]
}

func (c *closure) String() string {
	return "#<closure>"
}
