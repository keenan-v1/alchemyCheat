package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var combinations = make(map[string]*Combination)

// Combination x
type Combination struct {
	Part1         *Element
	Part2         *Element
	Result        *Element
	Completed     bool
	CompletedTime time.Time
}

func (c *Combination) String() string {
	complStr := ""
	if c.Completed {
		complStr = fmt.Sprintf(" [Completed: %s]", c.CompletedTime)
	}
	return fmt.Sprintf("%s = %s + %s\t\t%s", c.Result.Name, c.Part1.Name, c.Part2.Name, complStr)
}

// Key x
func (c *Combination) Key() string {
	return fmt.Sprintf("%d|%d", c.Part1.ID, c.Part2.ID)
}

// FindCombinationByIDs x
func FindCombinationByIDs(one int, two int) *Combination {
	if one < two {
		return combinations[fmt.Sprintf("%d|%d", one, two)]
	}
	return combinations[fmt.Sprintf("%d|%d", two, one)]
}

// FindCombinationByParts x
func FindCombinationByParts(one *Element, two *Element) *Combination {
	if one.ID < two.ID {
		return combinations[fmt.Sprintf("%d|%d", one.ID, two.ID)]
	}
	return combinations[fmt.Sprintf("%d|%d", two.ID, one.ID)]
}

// FindCombinationByKey x
func FindCombinationByKey(key string) *Combination {
	return combinations[key]
}

// NewCombination x
func NewCombination(one *Element, two *Element, result *Element) *Combination {
	c := &Combination{}
	if one.ID < two.ID {
		c.Part1 = one
		c.Part2 = two
	} else {
		c.Part1 = two
		c.Part2 = one
	}
	c.Result = result
	combinations[c.Key()] = c
	return c
}

// NewOrExistingCombination x
func NewOrExistingCombination(one *Element, two *Element, result *Element) *Combination {
	if combo := FindCombinationByParts(one, two); combo != nil {
		return combo
	}
	return NewCombination(one, two, result)
}

// ProcessHistory x
func ProcessHistory() {
	js, err := ioutil.ReadFile("history.json")
	check(err)

	// [{ "key": "12|345", "when": 123456789, "v": 123456789 }, ... ]
	data := []struct {
		Key  string `json:"key"`
		When int64  `json:"when"`
		V    int64  `json:"v"`
	}{}
	d := json.NewDecoder(bytes.NewReader(js))
	err = d.Decode(&data)
	check(err)
	for _, val := range data {
		combo := FindCombinationByKey(val.Key)
		if combo == nil {
			log.Printf("combo not found for %s - %s", val.Key, time.Unix(val.When, 0))
			continue
		}
		combo.Completed = true
		combo.CompletedTime = time.Unix(val.When, 0)
	}
}

// PrintCombinations x
func PrintCombinations() {
	line := 0
	for _, combo := range combinations {
		if combo.Completed || !combo.Part1.IsDiscovered() || !combo.Part2.IsDiscovered() {
			continue
		}
		line++
		fmt.Printf("#%4d %s [New: %t]\n", line, combo, !combo.Result.IsDiscovered())
	}
}

// DebugCombinations x
func DebugCombinations() {
	line := 0
	for key, combo := range combinations {
		if !combo.Completed {
			continue
		}
		line++
		fmt.Printf("#%4d (%7s) %s\n", line, key, combo)
	}
}
