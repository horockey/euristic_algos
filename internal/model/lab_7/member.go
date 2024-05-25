package model

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Member struct {
	Name     string
	Genotype Path
}

func NewMember(name string, g *Graph, s0 int) (*Member, error) {
	visited := map[int]struct{}{}
	p := make(Path, 0, len(g.Weights))
	currentNode := s0
	for len(visited) < len(g.Weights) {
		next := -1
		for _, found := visited[next]; next == -1 || found || next == currentNode; _, found = visited[next] {
			next = rnd.Intn(len(g.Weights))
			if len(visited) == len(g.Weights)-1 {
				goto ADD_START_TO_END
			}
		}

		p = append(p, &Transition{
			From:   currentNode,
			To:     next,
			Weight: g.Weights[currentNode][next],
		})
		visited[currentNode] = struct{}{}
		currentNode = next
	}

ADD_START_TO_END:
	p = append(p, &Transition{
		From:   currentNode,
		To:     s0,
		Weight: g.Weights[currentNode][s0],
	})

	return &Member{
		Name:     name,
		Genotype: p,
	}, nil
}

func (m *Member) String() string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", m.Name, m.Genotype.String())
}

func (m *Member) TwoPointMutation(i1, i2 int, g *Graph) (*Member, error) {
	if i1 < 0 || i1 >= len(m.Genotype) {
		return nil, fmt.Errorf("invalid i1: %d", i1)
	}
	if i2 < 0 || i2 >= len(m.Genotype) {
		return nil, fmt.Errorf("invalid i2: %d", i2)
	}

	if i2 < i1 {
		i1, i2 = i2, i1
	}

	if i1 == 0 {
		i1 = 1
	}
	if i2 == len(g.Weights)-1 {
		i2 = len(g.Weights) - 2
	}

	if i2 < i1 {
		i1, i2 = i2, i1
	}

	if i2 == i1 {
		if i1-1 > 0 {
			i1--
		} else if i2+1 < len(g.Weights) {
			i2++
		}
	}

	res := make(Path, 0, len(m.Genotype))
	for _, gen := range m.Genotype {
		res = append(res, &Transition{
			From:   gen.From,
			To:     gen.To,
			Weight: gen.Weight,
		})
	}

	res[i1] = &Transition{
		From:   res[i2].From,
		To:     res[i1].To,
		Weight: g.Weights[res[i2].From][res[i1].To],
	}
	if i1-1 >= 0 {
		res[i1-1].To = res[i1].From
	}

	res[i2] = &Transition{
		From:   m.Genotype[i1].From,
		To:     res[i2].To,
		Weight: g.Weights[m.Genotype[i1].From][res[i2].To],
	}
	if i2-1 >= 0 {
		res[i2-1].To = res[i2].From
	}

	return &Member{
		Name:     fmt.Sprintf("%s_mut%dx%d", m.Name, i1, i2),
		Genotype: res,
	}, nil
}

func (m *Member) TwoPointCrossover(other *Member, i1, i2 int, g *Graph, s0 int) (l, r *Member, resErr error) {
	if i2 < i1 {
		i1, i2 = i2, i1
	}

	if i1 == 0 {
		i1 = 1
	}
	if i2 == len(g.Weights)-1 {
		i2 = len(g.Weights) - 2
	}

	if i2 < i1 {
		i1, i2 = i2, i1
	}

	if i2 == i1 {
		if i1-1 > 0 {
			i1--
		} else if i2+1 < len(g.Weights) {
			i2++
		}
	}

	lNodes := map[int]struct{}{}
	for i := i1; i <= i2; i++ {
		lNodes[m.Genotype[i].To] = struct{}{}
	}
	rNodes := map[int]struct{}{}
	for i := i1; i <= i2; i++ {
		rNodes[other.Genotype[i].To] = struct{}{}
	}

	if !reflect.DeepEqual(lNodes, rNodes) {
		return m, other, nil
	}

	lGen := make(Path, 0, len(m.Genotype))
	rGen := make(Path, 0, len(m.Genotype))

	for i := 0; i < i1; i++ {
		lGen = append(lGen, m.Genotype[i])
		rGen = append(rGen, other.Genotype[i])
	}
	for i := i1; i < i2; i++ {
		rGen = append(rGen, &Transition{
			From:   rGen[len(rGen)-1].To,
			To:     m.Genotype[i].To,
			Weight: g.Weights[rGen[len(rGen)-1].To][m.Genotype[i].To],
		})
		lGen = append(lGen, &Transition{
			From:   lGen[len(lGen)-1].To,
			To:     other.Genotype[i].To,
			Weight: g.Weights[lGen[len(lGen)-1].To][other.Genotype[i].To],
		})
	}
	for i := i2; i < len(m.Genotype); i++ {
		lGen = append(lGen, &Transition{
			From:   lGen[len(lGen)-1].To,
			To:     m.Genotype[i].To,
			Weight: g.Weights[lGen[len(lGen)-1].To][m.Genotype[i].To],
		})
		rGen = append(rGen, &Transition{
			From:   rGen[len(rGen)-1].To,
			To:     other.Genotype[i].To,
			Weight: g.Weights[rGen[len(rGen)-1].To][other.Genotype[i].To],
		})
	}

	return &Member{
			Name:     fmt.Sprintf("%sx%s", m.Name, other.Name),
			Genotype: lGen,
		},
		&Member{
			Name:     fmt.Sprintf("%sx%s", other.Name, m.Name),
			Genotype: rGen,
		},
		nil
}
