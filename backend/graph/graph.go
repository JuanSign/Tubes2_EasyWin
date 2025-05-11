package graph

import (
	"fmt"
)

// json structure
type Element struct {
	Name   string     `json:"name"`
	Recipe [][]string `json:"recipe"`
}

type NodeType int

const (
	ElementNode NodeType = 0
	RecipeNode  NodeType = 1
)

type Node struct {
	Name string
	Type NodeType
}

type Graph struct {
	Nodes       []Node
	AdjList     map[int][]int
	NameToIndex map[string]int
}

// constructor
func NewGraph() *Graph {
	return &Graph{
		Nodes:       []Node{},
		AdjList:     make(map[int][]int),
		NameToIndex: make(map[string]int),
	}
}

// add node and return the index
func (g *Graph) addNode(name string, nodeType NodeType) int {
	if idx, exists := g.NameToIndex[name]; exists {
		return idx
	}
	idx := len(g.Nodes)
	g.Nodes = append(g.Nodes, Node{Name: name, Type: nodeType})
	g.NameToIndex[name] = idx
	return idx
}

// add adge
func (g *Graph) addEdge(src, dst int) {
	g.AdjList[src] = append(g.AdjList[src], dst)
}

// build graph from elements.json
func (g *Graph) BuildFromElements(elements []Element) {
	for _, el := range elements {
		elementIdx := g.addNode(el.Name, ElementNode)

		for _, recipe := range el.Recipe {
			if len(recipe) != 2 {
				continue
			}

			input1 := recipe[0]
			input2 := recipe[1]

			recipeName := fmt.Sprintf("%s + %s", input1, input2)
			recipeIdx := g.addNode(recipeName, RecipeNode)

			g.addEdge(elementIdx, recipeIdx)

			in1Idx := g.addNode(input1, ElementNode)
			in2Idx := g.addNode(input2, ElementNode)

			g.addEdge(recipeIdx, in1Idx)
			g.addEdge(recipeIdx, in2Idx)
		}
	}
}

// helper function to debug
func (g *Graph) DebugPrint() {

	// Print all nodes and their adjacency list
	for i, node := range g.Nodes {
		fmt.Printf("[%d] %s (%s) → ", i, node.Name, g.nodeTypeToString(node.Type))
		for _, adj := range g.AdjList[i] {
			fmt.Printf("%s (%s), ", g.Nodes[adj].Name, g.nodeTypeToString(g.Nodes[adj].Type))
		}
		fmt.Println()
	}

	// Print the number of nodes and edges
	fmt.Printf("Number of Nodes: %d\n", len(g.Nodes))
	edgesCount := 0
	for _, adjList := range g.AdjList {
		edgesCount += len(adjList)
	}
	fmt.Printf("Number of Edges: %d\n", edgesCount)
}

func (g *Graph) nodeTypeToString(t NodeType) string {
	if t == ElementNode {
		return "Element"
	}
	return "Recipe"
}
