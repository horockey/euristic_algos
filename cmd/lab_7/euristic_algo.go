package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	model "github.com/horockey/euristic_algos/internal/model/lab_7"
)

func euristicAlgo(
	g *model.Graph,
	s0 int,
	membersCount int,
	memberSize int,
	repeatsToBreak int,
	crossoverProba float32,
	mutationProba float32,
) (solution model.Path, resErr error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	bestResultRepeatedCount := 0
	lastBestResult := &model.Member{}

	g0, err := model.NewGeneration("g0", membersCount, g, s0)
	if err != nil {
		return nil, fmt.Errorf("creating generation: %w", err)
	}
	fmt.Printf("%s\nBest member: %s\n\n", g0.String(), g0.BestMember().String())

	for iteration := 0; bestResultRepeatedCount < repeatsToBreak; iteration++ {

		buffer := make([]*model.Member, 0, membersCount*4)
		buffer = append(buffer, g0.Members...)

		for _, member := range g0.Members {
			isCrosover := rnd.Float32() < crossoverProba
			isMutation := rnd.Float32() < mutationProba

			if isCrosover {
				var partner *model.Member
				for partner == nil || partner == member {
					partner = g0.Members[rnd.Intn(len(g0.Members))]
				}

				i1 := rand.Intn(memberSize)
				i2 := -1
				for i2 == -1 || i2 == i1 {
					i2 = rand.Intn(memberSize)
				}

				l, r, err := member.TwoPointCrossover(partner, i1, i2, g, s0)
				if err != nil {
					return nil, fmt.Errorf("crossingover: %w", err)
				}

				buffer = append(buffer, l, r)
			}

			if isMutation {
				i1 := rand.Intn(memberSize)
				i2 := -1
				for i2 == -1 || i2 == i1 {
					i2 = rand.Intn(memberSize)
				}

				mut, err := member.TwoPointMutation(i1, i2, g)
				if err != nil {
					return nil, fmt.Errorf("mutating: %w", err)
				}

				buffer = append(buffer, mut)
			}
		}

		uniqNames := map[string]struct{}{}
		buffer = slices.DeleteFunc(buffer, func(member *model.Member) bool {
			if _, found := uniqNames[member.Name]; found {
				return true
			}
			uniqNames[member.Name] = struct{}{}

			visitedTo := map[int]struct{}{}
			visitedFrom := map[int]struct{}{}

			for _, tr := range member.Genotype {
				if _, found := visitedTo[tr.To]; found {
					return true
				}
				visitedTo[tr.To] = struct{}{}

				if _, found := visitedFrom[tr.From]; found {
					return true
				}
				visitedFrom[tr.From] = struct{}{}

				if tr.Weight == 0 {
					return true
				}

				if tr.From == tr.To {
					return true
				}
			}
			return false
		})

		slices.SortFunc(buffer, func(a, b *model.Member) int {
			return a.Genotype.Total() - b.Genotype.Total()
		})

		if len(buffer) < membersCount {
			panic("buffer is to small")
		}

		gn := model.Generation{
			Name:    "g" + fmt.Sprint(iteration+1),
			Members: buffer[:membersCount],
		}

		fmt.Printf("%s\nBest member: %s\n\n", gn.String(), gn.BestMember().String())

		if res := gn.BestMember(); res.Genotype.Total() == lastBestResult.Genotype.Total() {
			bestResultRepeatedCount++
		} else {
			bestResultRepeatedCount = 1
			lastBestResult = res
		}

		g0 = &gn
	}

	return lastBestResult.Genotype, nil
}
