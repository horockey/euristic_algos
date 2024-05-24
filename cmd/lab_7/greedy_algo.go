package main

import (
	"fmt"

	"github.com/horockey/euristic_algos/internal/model"
)

func greedyAlgo(g *model.Graph, s0 int) (solution model.Path, resErr error) {
	visited := map[int]struct{}{}

	solution = make(model.Path, 0, len(g.Weights))
	currentNode := s0

	for len(visited) < len(g.Weights) {
		minIdx := -1
		for i, w := range g.Weights[currentNode] {
			if i == currentNode {
				continue
			}

			if w <= 0 && currentNode != i {
				return nil, fmt.Errorf(
					"invalid graph: %d->%d: %d",
					currentNode,
					i,
					w,
				)
			}

			if _, found := visited[i]; !found &&
				(minIdx == -1 || w < g.Weights[currentNode][minIdx]) {
				minIdx = i
			}
		}

		if minIdx == -1 {
			break
		}

		solution = append(solution, &model.Transition{
			From:   currentNode,
			To:     minIdx,
			Weight: g.Weights[currentNode][minIdx],
		})
		visited[currentNode] = struct{}{}
		currentNode = minIdx
	}

	solution = append(solution, &model.Transition{
		From:   currentNode,
		To:     s0,
		Weight: g.Weights[currentNode][s0],
	})
	return
}
