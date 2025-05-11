package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Element struct {
	Name   string     `json:"name"`
	Recipe [][]string `json:"recipe"`
}

func ScrapeElements(url string) ([]Element, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var elements []Element

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tbody := findChild(n, "tbody")
			if tbody != nil {
				trs := getChildren(tbody, "tr")
				if len(trs) > 0 {
					thNodes := getChildren(trs[0], "th")
					if len(thNodes) == 2 &&
						strings.TrimSpace(thNodes[0].FirstChild.Data) == "Element" &&
						strings.TrimSpace(thNodes[1].FirstChild.Data) == "Recipes" {
						for _, tr := range trs[1:] {
							tds := getChildren(tr, "td")
							if len(tds) == 2 {
								name := findFirstAnchor(tds[0])
								if name == "" {
									continue
								}
								var recipe [][]string
								ul := findChild(tds[1], "ul")
								if ul != nil {
									liNodes := getChildren(ul, "li")
									for _, li := range liNodes {
										recipePair := findRecipePair(li)
										if len(recipePair) == 2 {
											recipe = append(recipe, recipePair)
										}
									}
								}
								elements = append(elements, Element{Name: name, Recipe: recipe})
							}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return elements, nil
}

func getChildren(n *html.Node, tag string) []*html.Node {
	var result []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			result = append(result, c)
		}
	}
	return result
}

func findChild(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func findFirstAnchor(n *html.Node) string {
	var result string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					result = strings.TrimSpace(c.Data)
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result == "" {
				f(c)
			}
		}
	}
	f(n)
	return result
}

func findRecipePair(li *html.Node) []string {
	var recipePair []string
	anchorTags := getChildren(li, "a")
	if len(anchorTags) >= 2 {
		firstRecipe := findFirstAnchor(anchorTags[0])
		secondRecipe := findFirstAnchor(anchorTags[1])
		if firstRecipe != "" && secondRecipe != "" {
			recipePair = []string{firstRecipe, secondRecipe}
		}
	}
	return recipePair
}

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	elements, err := ScrapeElements(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of elements found: %d\n", len(elements))

	file, err := os.Create("../backend/elements.json")
	if err != nil {
		log.Fatal("Error creating JSON file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(elements)
	if err != nil {
		log.Fatal("Error encoding elements to JSON:", err)
	}

	fmt.Println("Scraped data has been saved to elements.json")
}
