package main

type ApplyFn func(s *Scope, args []Any) Any
type Function interface {
	Apply(*Scope, []Any) Any
}

type Macro interface {
	Function
	macro() bool
}

type closure struct {
	s          *Scope
	vars, body []Any
}

func (c *closure) Apply(s *Scope, args []Any) Any {
	args = c.s.EvalAll(args)
	c.s = NewScope(c.s)
	for i := 0; i < len(c.vars); i++ {
		v := c.vars[i]
		if v != Symbol("&") {
			c.s.env[v.(Symbol)] = args[i]
		} else {
			c.s.env[c.vars[i+1].(Symbol)] = args[i:]
			break
		}
	}

	evaled := c.s.EvalAll(c.body)
	return evaled[len(evaled)-1]
}

func (c *closure) String() string {
	return "#<closure>"
}

type macro closure

func (m *macro) Apply(s *Scope, args []Any) Any {
	m.s = NewScope(s)
	for i := 0; i < len(m.vars); i++ {
		v := m.vars[i]
		if v != Symbol("&") {
			m.s.env[v.(Symbol)] = args[i]
		} else {
			m.s.env[m.vars[i+1].(Symbol)] = args[i:]
			break
		}
	}

	m.s.EvalAll(m.body[:len(m.body)-1])

	return m.s.Eval(m.body[len(m.body)-1])
}

// Just a placeholder method that makes
// sure the macro's interface is unique
func (m *macro) macro() bool {
	return true
}

func (m *macro) String() string {
	return "#<macro>"
}
