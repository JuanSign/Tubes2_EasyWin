package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"backend/graph"
)

func loadElements(filePath string) ([]graph.Element, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var elements []graph.Element
	if err := json.NewDecoder(f).Decode(&elements); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return elements, nil
}

func main() {
	elements, err := loadElements("elements.json")
	if err != nil {
		log.Fatal(err)
	}

	g := graph.NewGraph()
	g.BuildFromElements(elements)

	// for debugging
	// g.DebugPrint()

	// example AllDFS usage
	// dfsResult := g.AllDFS("Air")

	// dfsResultJSON, err := json.MarshalIndent(dfsResult, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	return
	// }

	// err = os.WriteFile("dfs_result.json", dfsResultJSON, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing to file:", err)
	// 	return
	// }

	// fmt.Println("DFS result saved to dfs_result.json")
}
