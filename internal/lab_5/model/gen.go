package model

import "fmt"

type Gen struct {
	Load int
	Key  Key
}
type Key int

func (k Key) ProcNum(procCount int) int {
	pn := int(k) * procCount / 256
	if pn > procCount-1 {
		return procCount - 1
	}
	return pn
}

func NewGen(load int) *Gen {
	return &Gen{
		Load: load,
		Key:  Key(rnd.Intn(256)),
	}
}

func (g *Gen) Mutate() int {
	bit := rnd.Intn(8)
	g.Key = g.Key ^ (2 << (bit))
	return bit
}

func (g *Gen) String() string {
	return fmt.Sprintf("%d->%d", g.Load, g.Key)
}
