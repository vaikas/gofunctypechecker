package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/vaikas/buildpackstuff/pkg/detect"
)

func main() {
	// read all go files
	files, err := filepath.Glob("./testdata/*.go")
	if err != nil {
		log.Printf("failed to read directory %s\n", err)
		return
	}

	for _, f := range files {
		log.Printf("Processing file %s\n", f)
		// read file
		file, err := os.Open(f)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		// read the whole file in
		srcbuf, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}
		f := &detect.Function{File: f, Source: string(srcbuf)}
		if deets := detect.CheckFile(f); deets != nil {
			log.Printf("Found supported function in package %q signature %q", deets.Package, deets.Signature)
		}
	}
}