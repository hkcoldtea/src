package main

import (
	"image"
	"image/color"
	"log"

	"github.com/Nykakin/quantize"
	"github.com/oliamb/cutter"
)

func comparator(color1, color2 color.Color) float64 {
	const maxDiff = 765.0 // Difference between black and white colors

	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()

	r1, g1, b1 = r1>>8, g1>>8, b1>>8
	r2, g2, b2 = r2>>8, g2>>8, b2>>8

	return float64((max(r1, r2)-min(r1, r2))+
		(max(g1, g2)-min(g1, g2))+
		(max(b1, b2)-min(b1, b2))) / maxDiff
}

// min is minimum of two uint32
func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// max is maximum of two uint32
func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func GetDomainColor(img image.Image, op int) color.Color {
	rectangle := img.Bounds()
	sx := rectangle.Min.X
	sy := rectangle.Min.Y
	sw := rectangle.Max.X
	sh := rectangle.Max.Y

	switch op {
	case 0:
		sh = 10
	case 1:
		sy = sh - 10
		sh = 10
	case 2:
		sw = 10
	case 3:
		sx = sw - 10
		sw = 10
	}
	cImg, err := cutter.Crop(img, cutter.Config{
		Height:  sh,                  // height in pixel or Y ratio(see Ratio Option below)
		Width:   sw,                  // width in pixel or X ratio
		Mode:    cutter.TopLeft,      // Accepted Mode: TopLeft, Centered
		Anchor:  image.Point{sx, sy}, // Position of the top left point
		Options: cutter.Copy,
	})
	if err != nil {
		log.Fatal("Cannot crop image:", err)
	}

	quantizer := quantize.NewHierarhicalQuantizer()
	colors, err := quantizer.Quantize(cImg, 1)
	if err != nil {
		log.Fatal(err)
	}

	palette := make([]color.Color, len(colors))
	for index, clr := range colors {
		palette[index] = clr
	}

	if len(palette) == 0 {
		colorRGBA := color.RGBA{255, 255, 255, 255}
		return colorRGBA
	}

	return palette[0]
}

func GetLargeBorder(img image.Image, rx, ry int, tl, bl, ll, rl bool, threshold, ratio float64) image.Image {
	rectangle := img.Bounds()
	sx := rectangle.Min.X
	sy := rectangle.Min.Y
	sw := rectangle.Max.X
	sh := rectangle.Max.Y

	var colorRGBA color.Color
	match := 0

	if rx == -1 || ry == -1 {
		colorRGBA = GetDomainColor(img, 0)
	} else {
		r1, g1, b1, a1 := img.At(rx, ry).RGBA()
		r1, g1, b1, a1 = r1>>8, g1>>8, b1>>8, a1>>8
		colorRGBA = color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a1)}
	}

	if tl {
		var avgmatch float64

	TopLoop:
		for y := rectangle.Min.Y; y < rectangle.Max.Y-2; y++ {
			rectangle.Min.Y = y
			match = 0
			for x := rectangle.Min.X; x < rectangle.Max.X; x++ {
				if comparator(img.At(x, y), colorRGBA) < threshold {
					match++
				}
			}
			if avgmatch == 0.0 {
				avgmatch = float64(match)
			}
			avgmatch /= 2
			avgmatch += float64(match/2)
			if avgmatch / float64(rectangle.Max.X-rectangle.Min.X) >= ratio {
				sy++
			} else {
				break TopLoop
			}
		}
	}

	if bl {
		var avgmatch float64
		if rx == -1 || ry == -1 {
			colorRGBA = GetDomainColor(img, 1)
		}

	BottomLoop:
		for y := rectangle.Max.Y - 1; y >= rectangle.Min.Y; y-- {
			rectangle.Max.Y = y + 1
			match = 0
			for x := rectangle.Min.X; x < rectangle.Max.X; x++ {
				if comparator(img.At(x, y), colorRGBA) < threshold {
					match++
				}
			}
			if avgmatch == 0.0 {
				avgmatch = float64(match)
			}
			avgmatch /= 2
			avgmatch += float64(match/2)
			if avgmatch / float64(rectangle.Max.X-rectangle.Min.X) >= ratio {
				sh--
			} else {
				break BottomLoop
			}
		}
	}

	if ll {
		var avgmatch float64
		if rx == -1 || ry == -1 {
			colorRGBA = GetDomainColor(img, 2)
		}

	LeftLoop:
		for x := rectangle.Min.X; x < rectangle.Max.X-2; x++ {
			rectangle.Min.X = x
			match = 0
			for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
				if comparator(img.At(x, y), colorRGBA) < threshold {
					match++
				}
			}
			if avgmatch == 0.0 {
				avgmatch = float64(match)
			}
			avgmatch /= 2
			avgmatch += float64(match/2)
			if avgmatch / float64(rectangle.Max.X-rectangle.Min.X) >= ratio {
				sx++
			} else {
				break LeftLoop
			}
		}
	}

	if rl {
		var avgmatch float64
		if rx == -1 || ry == -1 {
			colorRGBA = GetDomainColor(img, 3)
		}

	RightLoop:
		for x := rectangle.Max.X - 1; x >= rectangle.Min.X; x-- {
			rectangle.Max.X = x + 1
			match = 0
			for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
				if comparator(img.At(x, y), colorRGBA) < threshold {
					match++
				}
			}
			if avgmatch == 0.0 {
				avgmatch = float64(match)
			}
			avgmatch /= 2
			avgmatch += float64(match/2)
			if avgmatch / float64(rectangle.Max.X-rectangle.Min.X) >= ratio {
				sw--
			} else {
				break RightLoop
			}
		}
	}

	croppedImage := image.Rectangle{image.Point{sx, sy}, image.Point{sw, sh}}
/*
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(croppedImage)
*/
	cropImg, _ := cropImage(img, croppedImage)
	return cropImg
}
