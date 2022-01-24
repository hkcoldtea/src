package main

import (
	"math"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

// ByHash implements sort.Interface based on the Hash field.
type ByHash []ItemStruct
/*
//Sort in reverse order
func (a ByHash) SortReverse() {
	sort.Sort(sort.Reverse(a))
}
*/
func (a ByHash) Len() int {
	return len(a)
}
func (a ByHash) Swap(i, j int) {
	(a)[i], (a)[j] = (a)[j], (a)[i]
}
func (a ByHash) Less(i, j int) bool {
	compI  := (a)[i].Compare
	compJ  := (a)[j].Compare
	if compI == 0 {
		(a)[i].PairName, (a)[i].Compare = a.MinCompare(i)
		compI = (a)[i].Compare
	}
	if compJ == 0 {
		(a)[j].PairName, (a)[j].Compare = a.MinCompare(j)
		compJ = (a)[j].Compare
	}

	compK  := hash_Compare((a)[i].Hash, (a)[j].Hash)
	if compJ > compK {
		(a)[j].PairName = (a)[i].Name
		(a)[j].Compare = compK
		compJ = compK
	}
	if compI > compK {
		(a)[i].PairName = (a)[j].Name
		(a)[i].Compare = compK
		compI = compK
	}

	return compI < compJ
}
func (a ByHash) MinCompare(j int) (string, float64) {
	var minIdx int
	var minValue float64 = math.MaxFloat64
	for i:=0; i<len(a); i++ {
		if i==j {
			continue
		}
		compK := hash_Compare((a)[i].Hash, (a)[j].Hash)
		if compK < minValue {
			minValue = compK
			minIdx = i
		}
	}
	return (a)[minIdx].Name, minValue
}

func hash_Compare(a, b *gocv.Mat) float64 {
	hash := contrib.ColorMomentHash{}
	return hash.Compare(*a, *b)
}
