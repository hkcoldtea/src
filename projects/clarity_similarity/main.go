package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

var (
	bSortByClarity = flag.Bool("sbc", false, "Sort by Clarity")
	bSortBySimilar = flag.Bool("sbs", false, "Sort by Similar")
	bVidDump       = flag.Bool("vcf", false, "Support VideoCaptureFile (debug)")
	bVerbose       = flag.Bool("v", false, "Verbose")
	bHelp          = flag.Bool("h", false, "print this help")
)

func main() {
	flag.Usage = func() {
		progname := filepath.Base(os.Args[0])
		fmt.Printf("How to run:\n  %s [-flags] [image.gif] [image.jpg] [video.mp4]\n\n", progname)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *bHelp {
		flag.Usage()
		return
	}
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No input filename given\n")
		flag.Usage()
		return
	}
	if *bSortBySimilar && *bSortByClarity {
		fmt.Fprintf(os.Stderr, "No supporting sort method\n")
		flag.Usage()
		return
	}
	if *bSortByClarity == false {
		*bSortBySimilar = true
	}

	var j int
	var filename string
	var items = make([]ItemStruct, flag.NArg()+1)

	for _, filename = range flag.Args() {
		fext := filepath.Ext(filename)
		if len(fext) > 0 {
			fext = strings.ToLower(fext[1:])
		}
		var f64Val float64
		var imghash *gocv.Mat
		switch fext {
		case "mp4":
			f64Val, imghash = handle_video(filename)
		case "gif":
			fallthrough
		case "jpeg":
			fallthrough
		case "jpg":
			fallthrough
		case "png":
			f64Val, imghash = handle_image(filename)
		default:
			f64Val = -1.0
		}

		if f64Val < 0 {
			// Remove the element at index i from items.
			copy(items[j:], items[j+1:]) // Shift left one index.
			items = items[:len(items)-1] // Truncate slice.
			continue
		}
		items[j].Name = filename
		items[j].Clarity = f64Val
		if *bSortBySimilar {
			items[j].Hash = imghash
		}
		j++
	}
	if j < 1 {
		return
	}
	items = items[:j]

	if *bSortByClarity {
		sort.Sort(ByClarity(items))
	} else {
		if *bSortBySimilar {
			sort.Sort(ByHash(items))
		}
	}
	fmt.Println("Order Clarity  Similar  Filename          Compare")
	for k, v := range items {
		k1 := k + 1
		if k == 0 {
			fmt.Printf("%3d   %6.3f %6.2f >> %q << %s\n", k1, v.Clarity, v.Compare, v.Name, v.PairName)
			continue
		}
		if k1 == flag.NArg() {
			fmt.Printf("%3d   %6.3f %6.2f -- %q -- %s\n", k1, v.Clarity, v.Compare, v.Name, v.PairName)
			continue
		}
		fmt.Printf("%3d   %6.3f %6.2f    %q    %s\n", k1, v.Clarity, v.Compare, v.Name, v.PairName)
	}
}

func handle_image(filename string) (float64, *gocv.Mat) {
	// read images
	mat := gocv.IMRead(filename, gocv.IMReadColor)
	defer mat.Close()
	if mat.Empty() {
		fmt.Fprintf(os.Stderr, "cannot read image %s\n", filename)
		return -1.0, nil
	}

	hash := contrib.ColorMomentHash{}
	resultB := mat.Clone()
//	defer resultB.Close()
	hash.Compute(mat, &resultB)

	val1, val2 := Clarity(mat)
	if *bVerbose {
		fmt.Println(filename, val1, val2)
	}

	if val2 > 20 {
		return val1, &resultB
	}
	return 0.0, &resultB
}

func handle_video(filename string) (float64, *gocv.Mat) {
	//load video
	vc, err := gocv.VideoCaptureFile(filename)
	if err != nil {
		return -1.0, nil
	}
	defer vc.Close()

	frames := vc.Get(gocv.VideoCaptureFrameCount)
	fps := vc.Get(gocv.VideoCaptureFPS)
	duration := frames / fps

	if *bVerbose {
		fmt.Println("frames=", frames, "fps=", fps, "duration=", duration)
	}

	sampletimeframe := 0.0
	for _, v := range []float64{64000, 800, 400, 20, 1} {
		if duration > v {
			sampletimeframe = v
			break
		}
	}
	timeframes := (sampletimeframe / duration) * frames

	mat := gocv.NewMat()
	defer mat.Close()
	matB := gocv.NewMat()
	defer matB.Close()

	var val1, val2 float64
	var maxValue, maxIdx float64

	resultA := gocv.NewMat()
//	defer resultA.Close()

	for {
		// Set Video frames
		vc.Set(gocv.VideoCapturePosFrames, timeframes)
		vc.Read(&mat)
		if mat.Empty() {
			break
		}

		hash := contrib.ColorMomentHash{}
		hash.Compute(mat, &resultA)

		val1, val2 = Clarity(mat)

		if *bVerbose {
			fmt.Println(filename, "timeframes=", timeframes, "val1=", val1, "val2=", val2)
		}

		if val2 > 20 {
			if val1 >= maxValue {
				maxValue = val1
				maxIdx = timeframes

				break
			}
		}

		timeframes += fps

		if timeframes > frames {
			break
		}
	}

	if maxIdx > 0.0 && *bVidDump {
		// Set Video frames
		vc.Set(gocv.VideoCapturePosFrames, maxIdx-1)
		vc.Read(&mat)
		vc.Read(&mat)
		if ! mat.Empty() {
			fname := filepath.Base(filename)
			fext := filepath.Ext(filename)
			fname = fname[0:len(fname)-len(fext)]
			saveFile := fmt.Sprintf("/tmp/clarity_%06.0f_%s.png", maxIdx, fname)
			gocv.IMWrite(saveFile, mat)
			val1, val2 = Clarity(mat)
			if *bVerbose {
				fmt.Println(saveFile, val1, val2)
			}
		}
	}

	return maxValue, &resultA
}
