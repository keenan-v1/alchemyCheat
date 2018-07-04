package types

import "sort"

// By x
type By func(e1, e2 *Element) bool

// Sort x
func (by By) Sort(elements []*Element) {
	s := &elementSorter{
		elements: elements,
		by:       by,
	}
	sort.Sort(s)
}

type elementSorter struct {
	elements []*Element
	by       func(e1, e2 *Element) bool
}

func (s *elementSorter) Len() int {
	return len(s.elements)
}

func (s *elementSorter) Swap(i, j int) {
	s.elements[i], s.elements[j] = s.elements[j], s.elements[i]
}

func (s *elementSorter) Less(i, j int) bool {
	return s.by(s.elements[i], s.elements[j])
}
