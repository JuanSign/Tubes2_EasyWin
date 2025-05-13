package graph

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

type Element struct {
	Name   string     `json:"name"`
	Recipe [][]string `json:"recipe"`
}

type NodeJSON struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Parent int    `json:"parent"`
}

type ReturnJSON struct {
	Name    string       `json:"name"`
	Content [][]NodeJSON `json:"content"`
}

type Node struct {
	Name string
}

type Graph struct {
	Nodes       []Node
	NameToIndex map[string]int
	Recipes     map[int][][2]int
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:       []Node{},
		NameToIndex: make(map[string]int),
		Recipes:     make(map[int][][2]int),
	}
}

func (g *Graph) addNode(name string) int {
	if idx, exists := g.NameToIndex[name]; exists {
		return idx
	}
	idx := len(g.Nodes)
	g.Nodes = append(g.Nodes, Node{Name: name})
	g.NameToIndex[name] = idx
	return idx
}

func (g *Graph) BuildFromElements(elements []Element) {
	for _, el := range elements {
		elementIdx := g.addNode(el.Name)
		for _, recipe := range el.Recipe {
			in1Idx := g.addNode(recipe[0])
			in2Idx := g.addNode(recipe[1])
			g.Recipes[elementIdx] = append(g.Recipes[elementIdx], [2]int{in1Idx, in2Idx})
		}
	}
}

func (g *Graph) AllDFS(start string) ReturnJSON {
	startIdx, exists := g.NameToIndex[start]
	if !exists {
		fmt.Println("Element not found!")
		return ReturnJSON{Name: start + " not found!"}
	}

	result := ReturnJSON{Name: start}
	var curID int
	visited := make(map[int]bool)

	var mu sync.Mutex
	var contentMu sync.Mutex
	var wg sync.WaitGroup

	var DFSTraversal func(idx int, parent int)
	DFSTraversal = func(idx int, parent int) {
		defer wg.Done()

		mu.Lock()
		if visited[idx] {
			mu.Unlock()
			return
		}
		visited[idx] = true
		mu.Unlock()

		for _, recipe := range g.Recipes[idx] {
			localContent := []NodeJSON{}

			mu.Lock()
			mergerID := curID + 1
			in1ID := curID + 2
			in2ID := curID + 3
			curID += 3
			mu.Unlock()

			localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: parent})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

			contentMu.Lock()
			result.Content = append(result.Content, localContent)
			contentMu.Unlock()

			wg.Add(2)
			go DFSTraversal(recipe[0], in1ID)
			go DFSTraversal(recipe[1], in2ID)
		}
	}

	wg.Add(1)
	go DFSTraversal(startIdx, 0)
	wg.Wait()

	return result
}

func (g *Graph) SingleDFS(start string) ReturnJSON {
	startIdx, exists := g.NameToIndex[start]
	if !exists {
		fmt.Println("Element not found!")
		return ReturnJSON{Name: start + " not found!"}
	}

	result := ReturnJSON{Name: start}
	var curID int
	visited := make(map[int]bool)

	var mu sync.Mutex
	var contentMu sync.Mutex
	var wg sync.WaitGroup

	var DFSTraversal func(idx int, parent int)
	DFSTraversal = func(idx int, parent int) {
		defer wg.Done()

		mu.Lock()
		if visited[idx] {
			mu.Unlock()
			return
		}
		visited[idx] = true
		mu.Unlock()

		recipes := g.Recipes[idx]
		if len(recipes) == 0 {
			return
		}

		randomIndex := rand.Intn(len(recipes))
		recipe := recipes[randomIndex]

		localContent := []NodeJSON{}

		mu.Lock()
		mergerID := curID + 1
		in1ID := curID + 2
		in2ID := curID + 3
		curID += 3
		mu.Unlock()

		localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: parent})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

		contentMu.Lock()
		result.Content = append(result.Content, localContent)
		contentMu.Unlock()

		wg.Add(2)
		go DFSTraversal(recipe[0], in1ID)
		go DFSTraversal(recipe[1], in2ID)
	}

	wg.Add(1)
	go DFSTraversal(startIdx, 0)
	wg.Wait()

	return result
}

func (g *Graph) AllBFS(start string) ReturnJSON {
	type QueueItem struct {
		Index  int
		Parent int
	}

	startIdx, exists := g.NameToIndex[start]
	if !exists {
		return ReturnJSON{Name: start + " not found!"}
	}

	result := ReturnJSON{Name: start}
	var curID int
	visited := make(map[int]bool)
	queue := []QueueItem{}

	for _, recipe := range g.Recipes[startIdx] {
		localContent := []NodeJSON{}

		mergerID := curID + 1
		in1ID := curID + 2
		in2ID := curID + 3
		curID += 3

		localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: 0})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

		result.Content = append(result.Content, localContent)

		queue = append(queue, QueueItem{Index: recipe[0], Parent: 0})
		queue = append(queue, QueueItem{Index: recipe[1], Parent: 0})
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Index] {
			continue
		}
		visited[current.Index] = true

		recipes := g.Recipes[current.Index]
		for _, recipe := range recipes {
			localContent := []NodeJSON{}

			mergerID := curID + 1
			in1ID := curID + 2
			in2ID := curID + 3
			curID += 3

			localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: current.Parent})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

			result.Content = append(result.Content, localContent)

			queue = append(queue, QueueItem{Index: recipe[0], Parent: in1ID})
			queue = append(queue, QueueItem{Index: recipe[1], Parent: in2ID})
		}
	}

	return result
}

func (g *Graph) SingleBFS(start string) ReturnJSON {
	type QueueItem struct {
		Index  int
		Parent int
		Root   int
	}

	startIdx, exists := g.NameToIndex[start]
	if !exists {
		return ReturnJSON{Name: start + " not found!"}
	}

	result := ReturnJSON{Name: start}
	visited := make(map[int]bool)
	queue := []QueueItem{}
	var curID int

	recipes := g.Recipes[startIdx]

	required := make([]int, len(recipes))
	for i := range required {
		required[i] = 2
	}

	for i, recipe := range recipes {
		localContent := []NodeJSON{}

		mergerID := curID + 1
		in1ID := curID + 2
		in2ID := curID + 3
		curID += 3

		localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: 0})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
		localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

		result.Content = append(result.Content, localContent)

		queue = append(queue, QueueItem{Index: recipe[0], Parent: in1ID, Root: i})
		queue = append(queue, QueueItem{Index: recipe[1], Parent: in2ID, Root: i})
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Index] {
			required[current.Root]--
			done := false
			for _, r := range required {
				if r <= 0 {
					done = true
					break
				}
			}
			if done {
				break
			}
			continue
		}
		visited[current.Index] = true

		for _, recipe := range g.Recipes[current.Index] {
			localContent := []NodeJSON{}

			mergerID := curID + 1
			in1ID := curID + 2
			in2ID := curID + 3
			curID += 3

			localContent = append(localContent, NodeJSON{Name: "merger", Id: mergerID, Parent: current.Parent})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: in1ID, Parent: mergerID})
			localContent = append(localContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: in2ID, Parent: mergerID})

			result.Content = append(result.Content, localContent)

			queue = append(queue, QueueItem{Index: recipe[0], Parent: in1ID, Root: current.Root})
			queue = append(queue, QueueItem{Index: recipe[1], Parent: in2ID, Root: current.Root})
		}
	}

	return result
}

func (g *Graph) DebugPrint() {
	for idx, node := range g.Nodes {
		fmt.Printf("[%d] %s → ", idx, node.Name)
		for _, pair := range g.Recipes[idx] {
			in1 := g.Nodes[pair[0]].Name
			in2 := g.Nodes[pair[1]].Name
			fmt.Printf("[%s + %s], ", in1, in2)
		}
		fmt.Println()
	}
	fmt.Printf("Number of Nodes: %d\n", len(g.Nodes))

	totalEdges := 0
	for _, recipeList := range g.Recipes {
		totalEdges += len(recipeList) * 2
	}
	fmt.Printf("Number of Edges: %d\n", totalEdges)
}

// API Handler
type RequestPayload struct {
	Element string `json:"element"`
	Type    string `json:"type"`
}

func (g *Graph) DFSHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var result ReturnJSON
	switch payload.Type {
	case "all":
		result = g.AllDFS(payload.Element)
	case "one":
		result = g.SingleDFS(payload.Element)
	default:
		http.Error(w, "Invalid type: must be 'one' or 'all'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (g *Graph) BFSHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var result ReturnJSON
	switch payload.Type {
	case "all":
		result = g.AllBFS(payload.Element)
	case "one":
		result = g.SingleBFS(payload.Element)
		return
	default:
		http.Error(w, "Invalid type: must be 'one' or 'all'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
