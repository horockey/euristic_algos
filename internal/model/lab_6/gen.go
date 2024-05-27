package model

import (
	"fmt"
	"math"
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

	for g.Task.Cost[g.ProcIdx()] == math.MaxInt {
		g.Key = Key(rnd.Intn(256))
	}

	return &g
}

func (g *Gen) SinglePointMutation() (bit int, changedProcNum bool) {
	oldProcNum := g.ProcIdx()

	for i := 0; i < 3; i++ {
		old := g.Key
		bit = rnd.Intn(8)
		g.Key = g.Key ^ (1 << (bit))

		if g.Task.Cost[g.ProcIdx()] < math.MaxInt {
			return bit, oldProcNum == g.ProcIdx()
		}
		g.Key = old
	}

	return -1, true
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
	cost := fmt.Sprintf("%d", g.Task.Cost[g.ProcIdx()])
	if g.Task.Cost[g.ProcIdx()] == math.MaxInt {
		cost = "∞"
	}
	return fmt.Sprintf("%s(%s)->%d", g.Task.Name, cost, g.Key)
}

func (g *Gen) Cost() int {
	return g.Task.Cost[g.ProcIdx()]
}
