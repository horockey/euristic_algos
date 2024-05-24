package model

import (
	"fmt"
)

type Gen struct {
	Task      *Task
	Key       Key
	procCount int
}
type Key int

func (g *Gen) ProcIdx() int {
	pn := int(g.Key) * g.procCount / 256
	if pn > g.procCount-1 {
		return g.procCount - 1
	}
	return pn
}

func NewGen(task *Task, procCount int) *Gen {
	g := Gen{
		Task:      task,
		Key:       Key(rnd.Intn(256)),
		procCount: procCount,
	}

	return &g
}

func (g *Gen) SinglePointMutation() (bit int, changedProcNum bool) {
	oldProcNum := g.ProcIdx()

	bit = rnd.Intn(8)
	g.Key = g.Key ^ (1 << (bit))

	return bit, oldProcNum == g.ProcIdx()
}

func (g *Gen) TwoPointMutation() (leftBit int, rightBit int, changedProcNum bool) {
	oldProcNum := g.ProcIdx()

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

	return leftBit, rightBit, oldProcNum == g.ProcIdx()
}

func (g *Gen) String() string {
	return fmt.Sprintf("%s(%d)->%d", g.Task.Name, g.Task.Cost[g.ProcIdx()], g.Key)
}

func (g *Gen) Cost() int {
	return g.Task.Cost[g.ProcIdx()]
}
