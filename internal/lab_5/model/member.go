package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Member struct {
	Genotype  []*Gen
	Name      string
	procCount int
}

func NewMember(loads []int, procCount int, name string) *Member {
	gt := make([]*Gen, 0, len(loads))
	for i := 0; i < len(loads); i++ {
		gt = append(gt, NewGen(loads[i]))
	}

	return &Member{
		Genotype:  gt,
		Name:      name,
		procCount: procCount,
	}
}

func (m *Member) Mutate() *Member {
	res := &Member{
		Name:      m.Name + "_mut",
		Genotype:  make([]*Gen, 0, len(m.Genotype)),
		procCount: m.procCount,
	}

	for _, gen := range m.Genotype {
		g := &Gen{
			Load: gen.Load,
			Key:  gen.Key,
		}
		res.Genotype = append(res.Genotype, g)
	}

	bit := res.Genotype[rnd.Intn(len(res.Genotype))].Mutate()
	res.Name += fmt.Sprint(bit)

	return res
}

func (m *Member) Crossover(partner *Member) (left *Member, right *Member) {
	left = &Member{
		Genotype:  make([]*Gen, len(m.Genotype)),
		Name:      fmt.Sprintf("(%sx%s)", m.Name, partner.Name),
		procCount: m.procCount,
	}
	right = &Member{
		Genotype:  make([]*Gen, len(m.Genotype)),
		Name:      fmt.Sprintf("(%sx%s)", partner.Name, m.Name),
		procCount: m.procCount,
	}

	split := rnd.Intn(len(m.Genotype))

	for i := 0; i < split; i++ {
		left.Genotype[i] = m.Genotype[i]
		right.Genotype[i] = partner.Genotype[i]
	}
	for i := split; i < len(m.Genotype); i++ {
		left.Genotype[i] = partner.Genotype[i]
		right.Genotype[i] = m.Genotype[i]
	}

	return
}

func (m *Member) Fenotype() Fenotype {
	ft := make([][]int, m.procCount)
	for i := 0; i < m.procCount; i++ {
		ft[i] = make([]int, 0, len(m.Genotype))
	}

	for i := 0; i < len(m.Genotype); i++ {
		ft[m.Genotype[i].Key.ProcNum(m.procCount)] = append(
			ft[m.Genotype[i].Key.ProcNum(m.procCount)],
			m.Genotype[i].Load,
		)
	}

	return ft
}

func (m *Member) String() string {
	res := m.Name + ": " + strings.Join(
		lo.Map(m.Genotype, func(item *Gen, _ int) string {
			return item.String()
		}),
		", ",
	)
	return res
}
