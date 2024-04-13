package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Generation struct {
	Members []*Member
	Name    string
}

func NewGeneration(
	membersCount int,
	memberSize int,
	procCount int,
	T1, T2 int,
	genName string,
) *Generation {
	loads := make([]int, 0, memberSize)
	for i := 0; i < memberSize; i++ {
		loads = append(loads, rnd.Intn(T2-T1)+T1)
	}

	mems := make([]*Member, 0, membersCount)
	for i := 0; i < membersCount; i++ {
		mems = append(mems,
			NewMember(
				loads,
				procCount,
				fmt.Sprintf("O%d", i),
			))
	}
	return &Generation{
		Members: mems,
		Name:    genName,
	}
}

func (g *Generation) BestMember() *Member {
	var bestMember *Member
	for _, member := range g.Members {
		if bestMember == nil ||
			member.Fenotype().Det() < bestMember.Fenotype().Det() {
			bestMember = member
		}
	}

	return bestMember
}

func (g *Generation) String() string {
	res := g.Name + ":\n\t" + strings.Join(lo.Map(g.Members, func(item *Member, _ int) string {
		return item.String()
	}), "\n\t")

	return res
}
