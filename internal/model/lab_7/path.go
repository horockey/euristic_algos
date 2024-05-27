package model

import (
	"fmt"
)

type Path []*Transition

type Transition struct {
	From   int
	To     int
	Weight int
}

func (p Path) String() string {
	if len(p) == 0 {
		return ""
	}

	res := fmt.Sprintf("%d", p[0].From)
	for i := 0; i < len(p); i++ {
		res += fmt.Sprintf(" -(%d)-> %d", p[i].Weight, p[i].To)
	}

	res += fmt.Sprintf(" (total: %d)", p.Total())

	return res
}

func (p Path) PumlString(color string) string {
	if len(p) == 0 {
		return ""
	}

	res := fmt.Sprintf("(%d) #gray\n", p[0].From)
	for _, tr := range p {
		res += fmt.Sprintf("(%d) -[#%s]-> (%d)\n",
			tr.From,
			color,
			tr.To,
		)
	}

	return res
}

func (p Path) Total() int {
	res := 0
	for _, tr := range p {
		res += tr.Weight
	}

	return res
}
