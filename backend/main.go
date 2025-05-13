package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	http.HandleFunc("/dfs", g.DFSHandler)
	http.HandleFunc("/bfs", g.BFSHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
