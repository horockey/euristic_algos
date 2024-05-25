package model

import (
	"fmt"
)

type Generation struct {
	Name    string
	Members []*Member
}

func NewGeneration(name string, membersCount int, g *Graph, s0 int) (*Generation, error) {
	gen := Generation{
		Name:    name,
		Members: make([]*Member, 0, membersCount),
	}

	for i := 0; i < membersCount; i++ {
		m, err := NewMember("O"+fmt.Sprint(i), g, s0)
		if err != nil {
			return nil, fmt.Errorf("creating member: %w", err)
		}

		gen.Members = append(gen.Members, m)
	}

	return &gen, nil
}

func (g *Generation) String() string {
	res := g.Name + "\n"
	for _, member := range g.Members {
		res += "\t" + member.String() + "\n"
	}

	return res
}

func (g *Generation) BestMember() *Member {
	res := g.Members[0]
	for i := 1; i < len(g.Members); i++ {
		if g.Members[i].Genotype.Total() < res.Genotype.Total() {
			res = g.Members[i]
		}
	}

	return res
}
