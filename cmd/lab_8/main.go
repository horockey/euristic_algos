package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"

	model "github.com/horockey/euristic_algos/internal/model/lab_6"
)

var tasks []*model.Task

func main() {
	var (
		procCount       int     = 5    // Кол-во процессоров
		tasksCount      int     = 11   // Кол-во задач
		loadMin         int     = 10   // Минимум нагрузки задачи
		loadMax         int     = 23   // Максимум нагрузки задачи
		repeatsToBreak  int     = 5    // Кол-во повторов для останова
		membersCount    int     = 5    // Кол-во особей в поколении
		crossoverProba  float32 = 1.0  // Вероятность кроссовера
		mutationProba   float32 = 1.0  // Вероятнотсь мутации
		impossibleProba float32 = 0.75 // Вероятность бесконечного времени исполнения
	)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Tasks:")
	fmt.Printf("%8s", "T/P")
	for i := 0; i < procCount; i++ {
		fmt.Printf("%5s", fmt.Sprintf("P%d", i+1))
	}
	fmt.Println()

	tasks = make([]*model.Task, 0, tasksCount)
	for i := 0; i < tasksCount; i++ {
		fmt.Printf("%5s\t", fmt.Sprintf("T%d", i+1))
		temp := make([]int, procCount)
		val := rand.Intn(loadMax-loadMin) + loadMin

		temp[rnd.Intn(len(temp))] = val

		for j := 0; j < procCount; j++ {
			if rnd.Float32() < impossibleProba && temp[j] == 0 {
				temp[j] = math.MaxInt
				fmt.Printf("%5s", "∞")
			} else {
				temp[j] = val
				fmt.Printf("%5d", temp[j])
			}
		}
		tasks = append(tasks, model.NewTask(fmt.Sprintf("T%d", i+1), temp))
		fmt.Println()
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
		g := model.NewGeneration(
			membersCount,
			tasks,
			fmt.Sprintf("gen_%d", k),
		)

		for i, member := range g0.Members {
			buffer := make([]*model.Member, 0, 4)
			buffer = append(buffer, member)

			isCrossover := rnd.Float32() < crossoverProba
			isMutation := rnd.Float32() < mutationProba

			if isCrossover {
				partnerIdx := rnd.Intn(len(g0.Members))
				if g0.Members[partnerIdx] == member {
					if partnerIdx < len(g0.Members)-1 {
						partnerIdx++
					} else {
						partnerIdx--
					}
				}

				l, r := member.SinglePointCrossover(g0.Members[partnerIdx])
				buffer = append(buffer, l)
				buffer = append(buffer, r)
			}

			if isMutation {
				buffer = append(buffer, member.SinglePointMutation())
			}

			slices.SortFunc(buffer, func(a *model.Member, b *model.Member) int {
				return a.Fenotype().Det() - b.Fenotype().Det()
			})

			g.Members[i] = buffer[0]
		}

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
