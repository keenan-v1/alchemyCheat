package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var elements = make(map[int]*Element)
var orderedElements = []*Element{}

// Condition x
type Condition struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}

func (c *Condition) String() string {
	return fmt.Sprintf("Type: %s, Total: %d", c.Type, c.Total)
}

// Element x
type Element struct {
	ID          int        `json:"id"`
	Prime       bool       `json:"prime"`
	Base        bool       `json:"base"`
	Name        string     `json:"name"`
	Hidden      bool       `json:"hidden"`
	Condition   *Condition `json:"condition"`
	ParentIDs   [][]int    `json:"parents"`
	ChildrenIDs []int      `json:"children"`
	Parents     []*Combination
	Children    []*Element
}

// IsDiscovered x
func (e *Element) IsDiscovered() bool {
	for _, p := range e.Parents {
		if p.Completed {
			return true
		}
	}
	return false
}

// IsCompleted x
func (e *Element) IsCompleted() bool {
	for _, p := range e.Parents {
		if !p.Completed {
			return false
		}
	}
	return true
}

// HasIncompleteChildren x
func (e *Element) HasIncompleteChildren() bool {
	for _, c := range e.Children {
		if !c.IsCompleted() && c.Combinations(e) != "" {
			return true
		}
	}
	return false
}

func (e *Element) String() string {
	return e.Name
}

// Combinations x
func (e *Element) Combinations(p *Element) string {
	matches := []string{}
	for _, combo := range e.Parents {
		if combo.Completed {
			continue
		}
		if combo.Part1.ID == p.ID {
			matches = append(matches, combo.Part2.Name)
			continue
		}
		if combo.Part2.ID == p.ID {
			matches = append(matches, combo.Part1.Name)
		}
	}
	if len(matches) > 0 {
		return strings.Join(matches, ", ")
	}
	return ""
}

// Info x
func (e *Element) Info() {
	fmt.Printf("Element: [%d] %s", e.ID, e.Name)
	if e.Prime {
		fmt.Printf(" [Prime]")
	}
	if e.Base {
		fmt.Printf(" [Base]")
	}
	if e.Hidden {
		fmt.Printf(" [Hidden]")
	}
	fmt.Printf("\n")
	if e.Condition != nil {
		fmt.Printf("Condition: %s\n", e.Condition)
	}

	if !e.IsCompleted() {
		fmt.Println("Parents:")
		for _, p := range e.Parents {
			if !p.Completed {
				fmt.Printf("\t%s\n", p)
			}
		}
	}
	if e.HasIncompleteChildren() {
		fmt.Println("Children:")
		for _, c := range e.Children {
			if !c.IsCompleted() {
				combos := c.Combinations(e)
				if len(combos) > 0 {
					fmt.Printf("\t%s = %s + %s\n", c.Name, e.Name, c.Combinations(e))
				}
			}
		}
	}
}

// AddElement x
func AddElement(e Element) {
	orderedElements = append(orderedElements, &e)
	elements[e.ID] = &e
}

// FindElementByID x
func FindElementByID(elementID int) *Element {
	return elements[elementID]
}

// ProcessElements x
func ProcessElements() {
	js, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic(err)
	}

	eArr := []Element{}
	err = json.Unmarshal(js, &eArr)
	if err != nil {
		panic(err)
	}

	for _, ele := range eArr {
		AddElement(ele)
	}
}

// PopulateCombinations x
func PopulateCombinations() {
	for _, e := range elements {
		for _, parents := range e.ParentIDs {
			combo := NewOrExistingCombination(FindElementByID(parents[0]), FindElementByID(parents[1]), e)
			e.Parents = append(e.Parents, combo)
		}
		if len(e.ChildrenIDs) > 0 {
			for _, childID := range e.ChildrenIDs {
				e.Children = append(e.Children, FindElementByID(childID))
			}
		}
	}
}

// PrintElements x
func PrintElements() {
	numberOfChildren := func(e1, e2 *Element) bool {
		return len(e1.Children) < len(e2.Children)
	}
	By(numberOfChildren).Sort(orderedElements)
	for _, e := range orderedElements {
		if !e.HasIncompleteChildren() {
			continue
		}
		e.Info()
		fmt.Println()
	}
}
