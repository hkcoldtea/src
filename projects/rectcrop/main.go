package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/webp"
)

var (
	BUILD string
)

func main() {
	input := flag.String("input", "", "input filename")
	output := flag.String("output", "", "output filename")
	outformat := flag.String("format", "", "output format")
	topEnabled := flag.Bool("top", false, "Top loop enabled")
	bottomEnabled := flag.Bool("bottom", false, "Bottom loop enabled")
	leftEnabled := flag.Bool("left", false, "Left loop enabled")
	rightEnabled := flag.Bool("right", false, "Right loop enabled")
	autoDisabled := flag.Bool("auto", false, "Disable auto-edge dectect")
	ref := flag.String("ref", "", "Reference points. e.g. 0,0")
	threshold := flag.Float64("threshold", 0.05, "Theshold")
	ratio := flag.Float64("ratio", 0.71, "Ratio")
	quality := flag.Int("quality", 85, "jpeg quality")
	binfo := flag.Bool("b", false, "Print build information")

	flag.Parse()

	if *binfo == true {
		fmt.Fprintf(os.Stderr, "Build: %s\n", BUILD)
		os.Exit(0)
	}

	if *input == "" {
		if len(flag.Args()) > 0 {
			*input = flag.Args()[0]
		}
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "No input file given")
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "No output file given")
		os.Exit(1)
	}

	img, format, err := readImage(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read image file: %v\n", err)
		os.Exit(1)
	}

	if *outformat != "" && *outformat != format {
		format = *outformat
	}

	out := *output

	var sx int
	var sy int
	if *ref == "" {
		sx = -1
		sy = -1
	} else {
		fmt.Sscanf(*ref, "%d,%d", &sx, &sy)
	}
	if *autoDisabled == false {
		*topEnabled = true
		*bottomEnabled = true
		*leftEnabled = true
		*rightEnabled = true
	}
	Newimg := GetLargeBorder(img, sx, sy, *topEnabled, *bottomEnabled, *leftEnabled, *rightEnabled, *threshold, *ratio)

	if float64(Newimg.Bounds().Size().X)*float64(Newimg.Bounds().Size().Y) <= 1 {
		Newimg, _ = cropCenter(img)
	}
	if err = writeImage(Newimg, out, format, *quality); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// readImage reads a image file from disk.
func readImage(name string) (image.Image, string, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, "", err
	}
	defer fd.Close()

	// image.Decode requires that you import the right image package.
	// We've imported "image/png", "image/gif", "image/jpeg".
	img, format, err := image.Decode(fd)
	if err != nil {
		if img, err = webp.Decode(fd); err != nil {
			return nil, "", err
		}
		format = "webp"
	}

	return img, format, nil
}

func cropCenter(img image.Image) (image.Image, error) {
	newX := float64(img.Bounds().Size().X) * 0.01
	newY := float64(img.Bounds().Size().Y) * 0.01

	// I've hard-coded a crop rectangle.
	return cropImage(img, image.Rect(int(newX*2), int(newY*2), int(newX*98), int(newY*98)))
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, out, format string, quality int) error {
	var fd io.WriteCloser
	var err error
	if out == "-" {
		fd = os.Stdout
	} else {
		fext := filepath.Ext(out)
		if len(fext) > 0 {
			fext = fext[1:]
		}
		if format != fext {
			if (format != "jpeg" || format != "jpg") && fext != "jpg" {
				out = out + "." + format
			}
		}
		fd, err = os.Create(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't create output file: %v\n", err)
			os.Exit(1)
		}
		defer fd.Close()
	}
	switch format {
	case "gif":
		var opt gif.Options
		opt.NumColors = 256
		return gif.Encode(fd, img, &opt)
	case "jpeg":
		fallthrough
	case "jpg":
		var opt jpeg.Options
		opt.Quality = quality
		return jpeg.Encode(fd, img, &opt)
	case "png":
		return png.Encode(fd, img)
	case "webp":
		err = errors.New("webp format does not support")
	default:
		fmt.Fprintln(os.Stderr, "unknown image format:", format)
	}

	return err
}
