package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/h2non/filetype"
)

// A result is the product of reading and summary.
type result struct {
	path    string
	summary string
	size    int64
}

var (
	build      *bool    = flag.Bool("b", false, "Prints project info")
	expectExt  *string  = flag.String("e", "", "Expected Extension, e.g.: gif")
	listonly   *bool    = flag.Bool("l", false, "Prints file list only")
	limSize    *int64   = flag.Int64("s", 0, "Minimum file size")
	threadNum  *int     = flag.Int("t", 2, "Number of thread")
	verbose    *bool    = flag.Bool("v", false, "Verbose")
	ignoreNew  *bool    = flag.Bool("3", false, "Ignore new-style compatibles")
	ignoreApp  *bool    = flag.Bool("A", false, "Ignore filetype = application/XXX")
	ignoreErr  *bool    = flag.Bool("E", false, "Ignore Error messages")
	ignoreFnt  *bool    = flag.Bool("F", false, "Ignore filetype = font/XXX")
	ignoreImg  *bool    = flag.Bool("I", false, "Ignore filetype = image/XXX")
	ignoreMed  *bool    = flag.Bool("M", false, "Ignore filetype = audio/XXX")
	ignoreTxt  *bool    = flag.Bool("T", false, "Ignore filetype = text/XXX")
	ignoreVid  *bool    = flag.Bool("V", false, "Ignore filetype = video/XXX")
	BUILD      string
)

func ReadBigFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		if *ignoreErr == false {
			fmt.Println("Can't opened this file")
		}
		return nil, err
    }
	defer f.Close()
	s := make([]byte, 512)
	for {
		f.Read(s[:])
		break
	}
	return s, nil
}

// I do work
func worker(id int, work chan result) {
	for i := range work {
		if i.size >= *limSize {
			if i.size == 0 {
				if *ignoreErr == false {
					fmt.Printf("skipped: open %s : empty\n", i.path)
				}
				continue
			}
			fext := filepath.Ext(i.path)
			if len(fext) > 0 {
				fext = fext[1:]
				if *expectExt != "" && strings.Compare(fext, *expectExt) != 0 {
					if *verbose == true {
						fmt.Printf("skipped: open %s : pattern no match\n", i.path)
					}
					continue
				}
			} else {
				if *expectExt != "." {
					if *ignoreErr == false {
						fmt.Printf("skipped: open %s : without file extension\n", i.path)
					}
					continue
				}
			}
			fext = strings.ToLower(fext)
			if *ignoreNew == false {
				switch fext {
				case "a":
					fext = "ar"
				case "aif":
					fext = "aiff"
				case "awk":
					fext = "txt"
				case "fla":
					fext = "flac"
				case "go":
					fext = "txt"
				case "htm":
					fext = "html"
				case "jpe":
					fallthrough
				case "jpeg":
					fext = "jpg"
				case "midi":
					fext = "mid"
				case "mpeg":
					fext = "mpg"
				case "mpeg4":
					fext = "mp4"
				case "pem":
					fext = "txt"
				case "sh":
					fext = "txt"
				case "sqlite3":
					fext = "sqlite"
				case "tiff":
					fext = "tif"
				case "text":
					fext = "txt"
				case "tgz":
					fext = "gz"
				case "py":
					fext = "txt"
				}
			}
			buf, err := ReadBigFile(i.path)
			if err != nil {
				if *ignoreErr == false {
					log.Println(err)
				}
				continue
			}
			kind, _ := filetype.Match(buf)
			if kind == filetype.Unknown || filetype.IsFont(buf) {
				contentType := http.DetectContentType(buf[:])
				// fallback
				if contentType == "application/octet-stream" {
					if *ignoreErr == false {
						fmt.Printf("skipped: open %s : unknown format\n", i.path)
					}
					continue
				} else {
					mime_extension, err := mime.ExtensionsByType(contentType)
					if err != nil {
						if *ignoreErr == false {
							fmt.Printf("skipped: open %s : unknown format\n", i.path)
						}
					} else {
						for _, v := range mime_extension {
							if len(v) > 0 {
								v = v[1:]
							}
							kind.Extension = v
							if fext == v {
								break
							}
						}
					}
					kind.MIME.Value = contentType
				}
			}
			if *ignoreApp && strings.Contains(kind.MIME.Value, "application/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore application format\n", i.path)
				}
				continue
			}
			if *ignoreMed && strings.Contains(kind.MIME.Value, "audio/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore audio format\n", i.path)
				}
				continue
			}
			if *ignoreFnt && strings.Contains(kind.MIME.Value, "font/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore font format\n", i.path)
				}
				continue
			}
			if *ignoreImg && strings.Contains(kind.MIME.Value, "image/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore image format\n", i.path)
				}
				continue
			}
			if *ignoreTxt && strings.Contains(kind.MIME.Value, "text/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore text format\n", i.path)
				}
				continue
			}
			if *ignoreVid && strings.Contains(kind.MIME.Value, "video/") {
				if *verbose == true {
					fmt.Printf("skipped: open %s : ignore video format\n", i.path)
				}
				continue
			}
			if strings.Compare(fext, kind.Extension) == 0 {
				if *verbose == true {
					fmt.Printf("matched: %q : MIME-type %q\n", i.path, kind.MIME.Value)
				}
				continue
			}
			if *listonly == false {
				fmt.Printf("filename: %q : suggested filetype: %q and MIME-type %q\n", i.path, kind.Extension, kind.MIME.Value)
			} else {
				fmt.Printf("%s\n", i.path)
			}
		}
	}
}

func main() {
	var totfiles int64
	var fpath []string

	flag.Parse()

	if *build {
		fmt.Println("Build info:", BUILD)
		os.Exit(0)
	}

	if *listonly == true {
		*ignoreErr = true
		if *verbose == true {
			fmt.Println("Can't use verbose and filename list at the same time.")
			return
		}
		*verbose = false
	}

	if *threadNum <= 0 {
		fmt.Println("Invalid number of threads")
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		fpath = flag.Args()
	} else {
		fpath = []string{"."}
	}

	wg := new(sync.WaitGroup)
	work := make(chan result)

	for i := 0; i < *threadNum; i++ {
		wg.Add(1)
		go func(i int) {
			worker(i, work)
			wg.Done()
		}(i)
	}

	fwalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if *ignoreErr == false {
				fmt.Println("skipped:", err)
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		var size int64
		size = info.Size()
		work <- result{path: path, size: size}
		totfiles += 1

		return nil
	}

	for _, cwd := range fpath {
		err := filepath.Walk(cwd, fwalk)
		if err != nil {
			log.Println(err)
		}
	}

	close(work)
	wg.Wait()

	if *listonly == false {
		fmt.Printf("%d files found\n", totfiles)
	}
}
