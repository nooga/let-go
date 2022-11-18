package vm

type Consts struct {
	consts []Value
	index  map[Value]int
}

func NewConsts() *Consts {
	return &Consts{
		consts: []Value{},
		index:  map[Value]int{},
	}
}

func (c *Consts) Intern(v Value) int {
	if i, ok := c.index[v]; ok {
		return i
	}
	c.consts = append(c.consts, v)
	i := len(c.consts) - 1
	c.index[v] = i
	return i
}

func (c *Consts) get(i int) Value {
	return c.consts[i]
}

func (c *Consts) count() int {
	return len(c.consts)
}
