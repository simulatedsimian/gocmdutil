package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func makeNewname(date time.Time, ext string, count int) string {
	return fmt.Sprintf("%04d_%02d_%02d_%02d_%02d_%02d_%02d%s",
		date.Year(), int(date.Month()), date.Day(),
		date.Hour(), date.Minute(), date.Second(), count, ext)
}

func main() {
	if len(os.Args) > 1 {
		for _, path := range os.Args[1:] {
			f, err := os.Open(path)
			exitOnError(err)
			defer f.Close()

			ex, err := exif.Decode(f)
			exitOnError(err)

			date, err := ex.DateTime()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error processing:", path, err)
				continue
			}

			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)

			count := 0
			for {

				newname := makeNewname(date, ext, count)
				newpath := filepath.Join(dir, newname)

				if fileExists(newpath) {
					count++
				} else {
					fmt.Printf("%s -> %s\n", path, newpath)
					os.Rename(path, newpath)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error processing:", path, err)
					}
					break
				}
			}

		}
	}
}
