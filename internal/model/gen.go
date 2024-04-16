package model

import (
	"fmt"
)

type Gen struct {
	Load      int
	Key       Key
	procCount int
}
type Key int

func (g *Gen) ProcNum() int {
	pn := int(g.Key) * g.procCount / 256
	if pn > g.procCount-1 {
		return g.procCount - 1
	}
	return pn
}

func NewGen(load int, procCount int) *Gen {
	return &Gen{
		Load:      load,
		Key:       Key(rnd.Intn(256)),
		procCount: procCount,
	}
}

func (g *Gen) SinglePointMutation() (bit int, changedProcNum bool) {
	oldProcNum := g.ProcNum()

	bit = rnd.Intn(8)
	g.Key = g.Key ^ (1 << (bit))

	return bit, oldProcNum == g.ProcNum()
}

func (g *Gen) TwoPointMutation() (leftBit int, rightBit int, changedProcNum bool) {
	oldProcNum := g.ProcNum()

	leftBit = rnd.Intn(8)
	rightBit = rnd.Intn(8)
	if leftBit < rightBit {
		leftBit, rightBit = rightBit, leftBit
	}
	if leftBit == rightBit {
		if rightBit < 7 {
			rightBit++
		} else if leftBit > 0 {
			leftBit++
		}
	}

	l := g.Key & (1 << leftBit)
	r := g.Key & (1 << rightBit)

	if l != r {
		g.Key = g.Key ^ (1 << leftBit)
		g.Key = g.Key ^ (1 << rightBit)
	}

	return leftBit, rightBit, oldProcNum == g.ProcNum()
}

func (g *Gen) String() string {
	return fmt.Sprintf("%d->%d", g.Load, g.Key)
}
