package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/horockey/euristic_algos/internal/model"
)

/*
	Модифицированный алгоритм Холланда.
	0. Матрица нагрузки неоднородная!
	1. Берем i-ю особь из поколения
	2. Определяем, будет ли двухточечный кроссовер (Pk)
		2.1. Если кроссовер - скрещиваем со случайной особью из текущего поколения
	3. Определяем, будет ли двухточечная мутация
		3.1. Если мутация - определяем ген и 2 бита в его ключе, меняем их местами
		3.2. Возможно будет смена держателя процесса
	4. i-й слот будущего поколения разыгрывается между Oi и его потомками (если есть)
*/

var tasks [][]int

func main() {
	var (
		procCount      int     = 5   // Кол-во процессоров
		tasksCount     int     = 11  // Кол-во задач
		loadMin        int     = 10  // Минимум нагрузки задачи
		loadMax        int     = 23  // Максимум нагрузки задачи
		repeatsToBreak int     = 5   // Кол-во повторов для останова
		membersCount   int     = 5   // Кол-во особей в поколении
		crossoverProba float32 = 1.0 // Вероятность кроссовера
		mutationProba  float32 = 1.0 // Вероятнотсь мутации
	)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	tasks = make([][]int, procCount)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = make([]int, tasksCount)
		for j := 0; j < len(tasks[i]); j++ {
			tasks[i][j] = rand.Intn(loadMax-loadMin) + loadMin
		}
	}

	fmt.Println()
	g0 := model.NewGeneration(
		membersCount,
		tasks,
		"gen_0",
	)

	fmt.Printf("%s\nBest fenotype (%s):\n%s\n",
		g0.String(),
		g0.BestMember().Name,
		g0.BestMember().Fenotype().String(),
	)

	lastBestResult := g0.BestMember().Fenotype().Det()
	bestResultRepeats := 1

	for k := 1; bestResultRepeats < repeatsToBreak; k++ {
		buffer := make([]*model.Member, 0, 2*membersCount)
		buffer = append(buffer, g0.Members...)

		for _, member := range g0.Members {
			isCrossover := rnd.Float32() < crossoverProba
			isMutation := rnd.Float32() < mutationProba

			candidates := make([]*model.Member, 0, 4)
			candidates = append(candidates, member)

			if isCrossover {
				partnerIdx := rnd.Intn(len(g0.Members))
				if g0.Members[partnerIdx] == member {
					if partnerIdx < len(g0.Members)-1 {
						partnerIdx++
					} else {
						partnerIdx--
					}
				}

				l, r := member.TwoPointCrossover(g0.Members[partnerIdx])
				candidates = append(candidates, l)
				candidates = append(candidates, r)
			}

			if isMutation {
				candidates = append(candidates, member.TwoPointMutaion())
			}

			var bestCandidate *model.Member
			for _, candidate := range candidates {
				if bestCandidate == nil ||
					candidate.Fenotype().Det() < bestCandidate.Fenotype().Det() {
					bestCandidate = candidate
				}
			}

			buffer = append(buffer, bestCandidate)
		}

		slices.SortFunc(buffer, func(a *model.Member, b *model.Member) int {
			return b.Fenotype().Det() - a.Fenotype().Det()
		})

		g := model.NewGeneration(
			membersCount,
			tasks,
			fmt.Sprintf("gen_%d", k),
		)
		// TODO: fill up

		fmt.Printf("%s\nBest fenotype (%s):\n%s\n",
			g.String(),
			g.BestMember().Name,
			g.BestMember().Fenotype().String(),
		)

		bestRes := g.BestMember().Fenotype().Det()
		if lastBestResult == bestRes {
			bestResultRepeats++
		} else {
			bestResultRepeats = 1
			lastBestResult = bestRes
		}

		g0 = g
	}
}
