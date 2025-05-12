package graph

import (
	"fmt"
	"strings"
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
		return ReturnJSON{}
	}

	result := ReturnJSON{Name: start}
	curID := 0
	visited := make(map[int]bool)

	var DFSTraversal func(idx int, parent int)
	DFSTraversal = func(idx int, parent int) {
		if visited[idx] {
			return
		}
		visited[idx] = true
		for _, recipe := range g.Recipes[idx] {
			curContent := []NodeJSON{}
			curContent = append(curContent, NodeJSON{Name: "merger", Id: curID + 1, Parent: parent})
			curContent = append(curContent, NodeJSON{Name: g.Nodes[recipe[0]].Name, Id: curID + 2, Parent: curID + 1})
			curContent = append(curContent, NodeJSON{Name: g.Nodes[recipe[1]].Name, Id: curID + 3, Parent: curID + 1})
			result.Content = append(result.Content, curContent)
			curID += 4
			DFSTraversal(recipe[0], curID-2)
			DFSTraversal(recipe[1], curID-1)
		}
	}

	DFSTraversal(startIdx, 0)
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
