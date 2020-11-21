package calc

type genericType int

const (
	gInt genericType = iota
	gFloat
)

type genericConst struct {
	value interface{}
	typ   genericType
}

func (g *genericConst) Add(val *genericConst) *genericConst {
	if g.typ == gFloat || val.typ == gFloat {
		return &genericConst{value: g.Float() + val.Float(), typ: gFloat}
	}
	return &genericConst{value: g.Int() + val.Int(), typ: gInt}
}

func (g *genericConst) Sub(val *genericConst) *genericConst {
	if g.typ == gFloat || val.typ == gFloat {
		return &genericConst{value: g.Float() - val.Float(), typ: gFloat}
	}
	return &genericConst{value: g.Int() - val.Int(), typ: gInt}
}

func (g *genericConst) Mul(val *genericConst) *genericConst {
	if g.typ == gFloat || val.typ == gFloat {
		return &genericConst{value: g.Float() * val.Float(), typ: gFloat}
	}
	return &genericConst{value: g.Int() * val.Int(), typ: gInt}
}

func (g *genericConst) Div(val *genericConst) *genericConst {
	return &genericConst{value: g.Float() / val.Float(), typ: gFloat}
}

func (g *genericConst) Float() float64 {
	if g.typ == gFloat {
		return g.value.(float64)
	}
	return float64(g.value.(int64))
}

func (g *genericConst) Int() int64 {
	if g.typ == gInt {
		return g.value.(int64)
	}
	return int64(g.value.(float64))
}

func (g *genericConst) Value() interface{} {
	return g.value
}
