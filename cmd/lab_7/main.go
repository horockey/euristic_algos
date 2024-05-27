package main

import (
	"errors"
	"fmt"
	"os"

	model "github.com/horockey/euristic_algos/internal/model/lab_7"
)

func main() {
	const (
		nodesCount     int     = 5   // Количество вершин графа
		weightMin      int     = 12  // Минимальный вес дуги
		weightMax      int     = 28  // Максимальный вес дуги
		membersCount   int     = 10  // Количество особей в поколении
		repeatsToBreak int     = 10  // Количетсво повторений для завершения
		crossoverProba float32 = 1.0 // Вероятность кроссовера
		mutationProba  float32 = 1.0 // Вероятность мутации
	)

	g := model.NewGraph(nodesCount, weightMin, weightMax)
	fmt.Println("Graph:")
	fmt.Println(g.String())

	fmt.Println("Greedy algo:")
	greedySol, err := greedyAlgo(g, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(greedySol.String())

	fmt.Println("Euristic algo:")
	euristicSol, err := euristicAlgo(
		g,
		0,
		membersCount,
		nodesCount,
		repeatsToBreak,
		crossoverProba,
		mutationProba,
	)
	if err != nil {
		panic(err)
	}

	if err := writeSolutionsToPuml("solution.puml", g, greedySol, euristicSol); err != nil {
		panic(err)
	}
}

func writeSolutionsToPuml(path string, g *model.Graph, greedyPath model.Path, eaPath model.Path) (resErr error) {
	fout, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating solution file: %w", err)
	}

	defer func() {
		if err := fout.Close(); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("closing file: %w", err))
		}
	}()

	res := "@startuml\n" +
		g.PumlString() +
		greedyPath.PumlString("red") +
		eaPath.PumlString("blue") +
		"@enduml\n"

	if _, err := fout.WriteString(res); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}
