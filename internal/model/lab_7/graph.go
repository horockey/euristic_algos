package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Graph struct {
	Weights [][]int
}

func NewGraph(nodesCount int, weightMin int, weightMax int) *Graph {
	g := Graph{
		Weights: make([][]int, nodesCount),
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < nodesCount; i++ {
		g.Weights[i] = make([]int, nodesCount)
		for j := 0; j < i; j++ {
			val := rnd.Intn(weightMax-weightMin) + weightMin
			g.Weights[i][j] = val
			g.Weights[j][i] = val
		}
	}

	return &g
}

func (g *Graph) String() string {
	res := ""
	for i := 0; i < len(g.Weights); i++ {
		for j := 0; j < len(g.Weights[i]); j++ {
			res += fmt.Sprintf("%5d", g.Weights[i][j])
		}
		res += "\n"
	}

	return res
}

func (g *Graph) PumlString() string {
	res := ""
	for i := 0; i < len(g.Weights); i++ {
		for j := 0; j < i; j++ {
			res += fmt.Sprintf("(%d) -- (%d): %d\n", i, j, g.Weights[i][j])
		}
	}

	return res
}
