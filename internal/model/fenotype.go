package model

import (
	"fmt"
	"math"
)

type Fenotype [][]int

func (f Fenotype) String() string {
	res := ""
	for row := 0; row < len(f); row++ {
		res += fmt.Sprintf("P%d: ", row+1)
		for col := 0; col < len(f[row]); col++ {
			res += fmt.Sprintf("%5d", f[row][col])
		}
		res += fmt.Sprintf(" (sum: %d)\n", sum(f[row]))
	}
	res += fmt.Sprintf("det: %d\n", f.Det())
	return res
}

func (f Fenotype) Det() int {
	m := math.MinInt
	for row := range f {
		s := sum(f[row])
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
