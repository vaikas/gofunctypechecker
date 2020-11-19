package main

import (
	"fmt"
	"github.com/vaikas/buildpackstuff/pkg/detect"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const plan = `
[[provides]]
name = "go-dist"

[[provides]]
name = "go-mod-vendor"

[[provides]]
name = "go-build"

[[requires]]
name = "go-dist"
version = "gcr.io/paketo-buildpacks/go-dist:0.2.5"

[[requires]]
name = "go-mod-vendor"
version = "gcr.io/paketo-buildpacks/go-mod-vendor:0.0.169"

[[requires]]
name = "go-build"
version = "gcr.io/paketo-buildpacks/go-build:0.1.1"
`

const supportedFuncs = `
import (
    "context"
	event "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

The following function signatures are supported by this builder:
func(event.Event)
func(event.Event) protocol.Result
func(context.Context, event.Event)
func(context.Context, event.Event) protocol.Result
func(event.Event) *event.Event
func(event.Event) (*event.Event, protocol.Result)
func(context.Context, event.Event) *event.Event
func(context.Context, event.Event) (*event.Event, protocol.Result)
`

func printSupportedFunctionsAndExit() {
	fmt.Println(supportedFuncs)
	os.Exit(100)
}

func main() {
	log.Println("ARGS: ", os.Args)
	for _, e := range os.Environ() {
		log.Println(e)
	}

	planFile, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("failed to open the plan file for writing", os.Args[2], err)
		printSupportedFunctionsAndExit()
	}
	defer planFile.Close()
	if _, err := planFile.WriteString(string(plan)); err != nil {
		printSupportedFunctionsAndExit()
	}

	// read all go files
	files, err := filepath.Glob("./*.go")
	if err != nil {
		log.Printf("failed to read directory %s\n", err)
		printSupportedFunctionsAndExit()
	}

	for _, f := range files {
		log.Printf("Processing file %s\n", f)
		// read file
		file, err := os.Open(f)
		if err != nil {
			log.Println(err)
			printSupportedFunctionsAndExit()
		}
		defer file.Close()

		// read the whole file in
		srcbuf, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			printSupportedFunctionsAndExit()
		}
		f := &detect.Function{File: f, Source: string(srcbuf)}
		if deets := detect.CheckFile(f); deets != nil {
			log.Printf("Found supported function in package %q signature %q", deets.Package, deets.Signature)
			os.Exit(0)
			if err := os.Mkdir("TESTDIR", os.ModeDir); err != nil {
				log.Println("FAILED TO MKDIR: ", err)
			}
		}
	}
	printSupportedFunctionsAndExit()
}