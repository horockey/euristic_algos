package model

import (
	"fmt"
	"math"

	"github.com/samber/lo"
)

type Fenotype [][]*Gen

func (f Fenotype) String() string {
	res := ""
	for row := 0; row < len(f); row++ {
		res += fmt.Sprintf("P%d: ", row+1)
		for col := 0; col < len(f[row]); col++ {
			cost := fmt.Sprintf("%d", f[row][col].Cost())
			if f[row][col].Cost() == math.MaxInt {
				cost = "∞"
			}
			res += fmt.Sprintf("%s(%s) ", f[row][col].Task.Name, cost)
		}
		res += fmt.Sprintf(" (sum: %d)\n", lo.SumBy(f[row], func(g *Gen) int { return g.Cost() }))
	}
	res += fmt.Sprintf("det: %d\n", f.Det())
	return res
}

func (f Fenotype) Det() int {
	m := math.MinInt
	for row := range f {
		s := lo.SumBy(f[row], func(g *Gen) int { return g.Cost() })
		if s > m {
			m = s
		}
	}
	return m
}

func sum(a []int) int {
	s := 0
	for _, val := range a {
		s += val
	}
	return s
}

func max(a []int) int {
	m := math.MinInt
	for _, val := range a {
		if val > m {
			m = val
		}
	}
	return m
}
