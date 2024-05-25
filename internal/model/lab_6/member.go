package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

type Member struct {
	Genotype  []*Gen
	Name      string
	procCount int
}

func NewMember(tasks []*Task, name string) *Member {
	gt := make([]*Gen, 0, len(tasks))
	ts := slices.Clone(tasks)
	for i := 0; i < len(tasks); i++ {
		idx := rnd.Intn(len(ts))
		gt = append(gt, NewGen(ts[idx], len(tasks[0].Cost)))
		ts = append(ts[:idx], ts[idx+1:]...)
	}

	return &Member{
		Genotype:  gt,
		Name:      name,
		procCount: len(tasks[0].Cost),
	}
}

func (m *Member) SinglePointMutation() *Member {
	res := &Member{
		Name:      m.Name + "_mut",
		Genotype:  make([]*Gen, 0, len(m.Genotype)),
		procCount: m.procCount,
	}

	for _, gen := range m.Genotype {
		g := &Gen{
			Task:      gen.Task,
			Key:       gen.Key,
			procCount: m.procCount,
		}
		res.Genotype = append(res.Genotype, g)
	}

	bit, _ := res.Genotype[rnd.Intn(len(res.Genotype))].SinglePointMutation()
	res.Name += fmt.Sprint(bit)

	return res
}

func (m *Member) TwoPointMutaion() *Member {
	res := &Member{
		Name:      m.Name + "_mut",
		Genotype:  make([]*Gen, 0, len(m.Genotype)),
		procCount: m.procCount,
	}

	for _, gen := range m.Genotype {
		g := &Gen{
			Task:      gen.Task,
			Key:       gen.Key,
			procCount: m.procCount,
		}
		res.Genotype = append(res.Genotype, g)
	}

	lBit, rBit, _ := res.Genotype[rnd.Intn(len(res.Genotype))].TwoPointMutation()
	res.Name += fmt.Sprintf("%d.%d", lBit, rBit)

	return res
}

func (m *Member) SinglePointCrossover(partner *Member) (left *Member, right *Member) {
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

func (m *Member) TwoPointCrossover(partner *Member) (left *Member, right *Member) {
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

	splitL := rnd.Intn(len(m.Genotype))
	splitR := rnd.Intn(len(m.Genotype))
	if splitL > splitR {
		splitL, splitR = splitR, splitL
	}
	if splitL == splitR {
		if splitR+1 < len(m.Genotype) {
			splitR++
		} else if splitL-1 > 0 {
			splitL--
		}
	}

	for i := 0; i < splitL; i++ {
		left.Genotype[i] = m.Genotype[i]
		right.Genotype[i] = partner.Genotype[i]
	}
	for i := splitL; i < splitR; i++ {
		left.Genotype[i] = partner.Genotype[i]
		right.Genotype[i] = m.Genotype[i]
	}
	for i := 0; i < len(m.Genotype); i++ {
		left.Genotype[i] = m.Genotype[i]
		right.Genotype[i] = partner.Genotype[i]
	}

	return left, right
}

func (m *Member) Fenotype() Fenotype {
	ft := make([][]*Gen, m.procCount)
	for i := 0; i < m.procCount; i++ {
		ft[i] = make([]*Gen, 0, len(m.Genotype))
	}

	for i := 0; i < len(m.Genotype); i++ {
		ft[m.Genotype[i].ProcIdx()] = append(
			ft[m.Genotype[i].ProcIdx()],
			m.Genotype[i],
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
