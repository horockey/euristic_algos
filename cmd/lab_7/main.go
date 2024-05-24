package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/horockey/euristic_algos/internal/model"
)

func main() {
	const (
		nodesCount int = 5  // Количество вершин графа
		weightMin  int = 12 // Минимальный вес дуги
		weightMax  int = 28 // Максимальный вес дуги
	)

	g := model.NewGraph(nodesCount, weightMin, weightMax)
	fmt.Println("Graph:")
	fmt.Println(g.String())

	fmt.Println("GreedyAlgo:")
	greedySol, err := greedyAlgo(g, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(greedySol.String())

	if err := writeSolutionsToPuml("solution.puml", g, greedySol, nil); err != nil {
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
