package model

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
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
	return fmt.Sprintf("%s: %s", m.Name, m.Genotype.String())
}

func (m *Member) TwoPointMutation(i1, i2 int, g *Graph) (*Member, error) {
	p := m.Genotype

	if i1 < 0 || i1 >= len(p) {
		return nil, fmt.Errorf("invalid i1: %d", i1)
	}
	if i2 < 0 || i2 >= len(p) {
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
		}
		if i2+1 < len(g.Weights) {
			i2++
		}
	}

	res := slices.Clone(p)

	res[i1] = &Transition{
		From:   res[i1].From,
		To:     res[i2].To,
		Weight: g.Weights[res[i1].From][res[i2].To],
	}
	if i1+1 < len(res) {
		res[i1+1].From = res[i1].To
	}

	res[i2] = &Transition{
		From:   res[i2].From,
		To:     m.Genotype[i1].To,
		Weight: g.Weights[res[i2].From][m.Genotype[i1].To],
	}
	if i2+1 < len(res) {
		res[i2+1].From = res[i2].To
	}

	return &Member{
		Name:     fmt.Sprintf("%s_mut%dx%d", m.Name, i1, i2),
		Genotype: res,
	}, nil
}

func (m *Member) TwoPointCrossover(other *Member, i1, i2 int, g *Graph, s0 int) (l, r *Member, resErr error) {
	p := m.Genotype

	lNodes := map[int]struct{}{}
	for _, el := range m.Genotype {
		lNodes[el.To] = struct{}{}
	}
	rNodes := map[int]struct{}{}
	for _, el := range other.Genotype {
		rNodes[el.To] = struct{}{}
	}

	if i1 <= 0 {
		i1 = 1
	}
	if i1 >= len(p)-1 {
		i1 = len(p) - 2
	}

	if i2 <= 0 {
		i2 = 1
	}
	if i2 >= len(p)-1 {
		i2 = len(p) - 2
	}

	lGen := make(Path, 0, len(p))
	rGen := make(Path, 0, len(p))

	delete(lNodes, s0)
	delete(rNodes, s0)

	for i := 0; i < i1; i++ {
		lGen = append(lGen, p[i])
		delete(lNodes, p[i].To)

		rGen = append(rGen, other.Genotype[i])
		delete(rNodes, other.Genotype[i].To)
	}
	for i := i1; i < i2; i++ {
		_, found := lNodes[other.Genotype[i].To]
		if !found {
			continue
		}
		delete(lNodes, other.Genotype[i].To)

		from := m.Genotype[0].From
		if len(lGen) > 0 {
			from = lGen[len(lGen)-1].To
		}
		to := other.Genotype[i].To
		lGen = append(lGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})

		_, found = rNodes[m.Genotype[i].To]
		if !found {
			continue
		}
		delete(rNodes, m.Genotype[i].To)

		from = other.Genotype[0].From
		if len(rGen) > 0 {
			from = rGen[len(rGen)-1].To
		}
		to = m.Genotype[i].To
		rGen = append(rGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})
	}
	for i := i2; i < len(p); i++ {
		_, found := rNodes[other.Genotype[i].To]
		if !found {
			continue
		}
		delete(rNodes, other.Genotype[i].To)

		from := rGen[len(rGen)-1].To
		to := other.Genotype[i].To
		rGen = append(rGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})

		_, found = lNodes[m.Genotype[i].To]
		if !found {
			continue
		}
		delete(lNodes, m.Genotype[i].To)

		from = lGen[len(lGen)-1].To
		to = m.Genotype[i].To
		lGen = append(lGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})
	}

	for node := range lNodes {
		if node == s0 {
			continue
		}

		from := lGen[len(lGen)-1].To
		to := node
		lGen = append(lGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})
	}

	if lGen[len(lGen)-1].To != s0 {
		lGen = append(lGen, &Transition{
			From:   lGen[len(lGen)-1].To,
			To:     s0,
			Weight: g.Weights[lGen[len(lGen)-1].To][s0],
		})
	}

	for node := range rNodes {
		if node == s0 {
			continue
		}

		from := rGen[len(rGen)-1].To
		to := node
		rGen = append(rGen, &Transition{
			From:   from,
			To:     to,
			Weight: g.Weights[from][to],
		})
	}

	if rGen[len(rGen)-1].To != s0 {
		rGen = append(rGen, &Transition{
			From:   rGen[len(rGen)-1].To,
			To:     s0,
			Weight: g.Weights[rGen[len(rGen)-1].To][s0],
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

func (m *Member) Crossover(other *Member, i1, i2 int, g *Graph, s0 int) (l, r *Member, resErr error) {
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
		}
		if i2+1 < len(g.Weights) {
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
			From:   lGen[len(rGen)-1].To,
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
