package main

import (
	"gocv.io/x/gocv"
)

//Custom structure, used to customize sorting
type ItemStruct struct {
	Name     string
	Clarity  float64
	Hash     *gocv.Mat
	PairName string
	Compare  float64
}
