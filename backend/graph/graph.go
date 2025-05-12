package graph

import (
	"fmt"
	"strings"
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
			curID += 4
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

func (g *Graph) BFS(start string) {
	startIdx, exists := g.NameToIndex[start]
	if !exists {
		fmt.Println("Element not found!")
		return
	}
	visited := make(map[int]bool)

	var BFS func(idx int, depth int)
	BFS = func(idx int, depth int) {

		queue := []int{}
		queue = append(queue, idx)

		for len(queue) > 0 {

			current := queue[0]
			queue = queue[1:]
			fmt.Printf("%s%s\n", strings.Repeat("-", depth), g.Nodes[current].Name)
			if visited[current] {
				continue
			}
			visited[current] = true
			recipes := g.Recipes[current]
			for _, recipe := range recipes {
				recipe1 := recipe[0]
				recipe2 := recipe[1]
				queue = append(queue, recipe1, recipe2)
			}
			depth++

		}

	}
	BFS(startIdx, 0)

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
