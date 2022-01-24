package main

import (
	"sort"
)

// ByClarity implements sort.Interface based on the Hash field.
type ByClarity []ItemStruct

//Sort in reverse order
func (a ByClarity) SortReverse() {
	sort.Sort(sort.Reverse(a))
}
func (a ByClarity) Len() int {
	return len(a)
}
func (a ByClarity) Swap(i, j int) {
	(a)[i], (a)[j] = (a)[j], (a)[i]
}

func (a ByClarity) Less(i, j int) bool {
	return (a)[i].Clarity < (a)[j].Clarity
}
