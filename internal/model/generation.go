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
	tasks [][]int,
	genName string,
) *Generation {
	loads := make([]int, 0, len(tasks[0]))
	for i := 0; i < len(tasks[0]); i++ {
		loads = append(loads, tasks[rnd.Intn(len(tasks))][i])
	}

	mems := make([]*Member, 0, membersCount)
	for i := 0; i < membersCount; i++ {
		mems = append(mems,
			NewMember(
				loads,
				len(tasks),
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
