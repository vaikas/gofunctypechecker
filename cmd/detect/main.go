package main

import (
	"bufio"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vaikas/gofunctypechecker/pkg/detect"
)

const supportedFuncs = `
Could not find a supported function signature. Supported function signatures
are listed below. Note that the function must be visible outside of the package
(capitalized, for example, Handler vs. handler).
Supported function Signatures:
%s
`

type PlanArguments struct {
	Name string
	Package string
	Function string
	Env map[string]string
}

const defaultPlanName = "http-go-function"

const defaultPlanTemplate = `
[[provides]]
name = "{{ .Name }}"
[[requires]]
name = "{{ .Name }}"
[requires.metadata]
package = "{{ .Package }}"
function = "{{ .Function }}"
exampleEnv = "{{ index .Env "PWD" }}"
`

// defaultSignatures that just use HTTP handler.
var defaultSignatures = []detect.FunctionSignature{
	{In: []detect.FunctionArg{
		{ImportPath: "net/http", Name: "ResponseWriter"},
		{ImportPath: "net/http", Name: "Request", Pointer: true},
	}}}

type EnvConfig struct {
	GoPackage  string `envconfig:"GO_PACKAGE" default:"./"`
	GoFunction string `envconfig:"GO_FUNCTION" default:"Receiver"`
	Protocol   string `envconfig:"PROTOCOL" default:"http"`
	Signatures string `envconfig:"SIGNATURES"`
	PlanTemplate string `envconfig:"PLAN_TEMPLATE"`
}

func printSupportedFunctionsAndExit(sigs string) {
	fmt.Println(supportedFuncs)
	os.Exit(100)
}

func main() {
	log.Println("ARGS: ", os.Args)
	for _, e := range os.Environ() {
		log.Println(e)
	}

	if len(os.Args) < 3 {
		log.Println("Usage: %s <PLATFORM_DIR> <BUILD_PLAN>", os.Args[0])
		os.Exit(100)
	}

	planFileName := os.Args[2]
	log.Println("using plan file: ", planFileName)

	// Grab the env variables
	var envConfig EnvConfig
	if err := envconfig.Process("detect", &envConfig); err != nil {
		log.Printf("Failed to process env variables: %s\n", err)
	}

	moduleName, err := readModuleName()
	if err != nil {
		log.Println("Failed to read go.mod file: ", err)
		os.Exit(100)
	}
	// There are two ENV variables that control what should be checked.
	// We yank the base package from go.mod and append CE_GO_PACKAGE into it
	// if it's given.
	// TODO: Use library for these...
	goPackage := envConfig.GoPackage
	if !strings.HasSuffix(goPackage, "/") {
		goPackage = goPackage + "/"
	}
	// fullGoPackage is the import path that we'll use with the scaffolding.
	// Looks for example like: github.com/vaikas/buildpackstuffhttp/pkg/detect/testdata
	fullGoPackage := moduleName
	if goPackage != "./" {
		fullGoPackage = fullGoPackage + "/" + filepath.Clean(goPackage)
	}
	log.Println("Using relative path to look for function: ", goPackage)

	goFunction := envConfig.GoFunction

	// Construct the detector. Either using default, fetch the config from a URL or from a file.
	var detector *detect.Detector
	if envConfig.Signatures == "" {
		// Just use defaults
		detector = detect.NewDetector(defaultSignatures)
	} else if (strings.HasSuffix(envConfig.Signatures, "http://") || strings.HasSuffix(envConfig.Signatures, "https://")) {
		detector, err = detect.NewDetectorFromURL(envConfig.Signatures)
		if err != nil {
			log.Fatalf("Failed to create detector with signatures from URL %q : %s\n", envConfig.Signatures, err)
			os.Exit(100)
		}
	} else {
		detector, err = detect.NewDetectorFromFile(envConfig.Signatures)
		if err != nil {
			log.Fatalf("Failed to create detector with signatures from file %q : %s\n", envConfig.Signatures, err)
			os.Exit(100)
		}
	}

	// Grab the plan template if specified, or use the default.
	planTemplate := defaultPlanTemplate
	if envConfig.PlanTemplate != "" {
		if (strings.HasSuffix(envConfig.PlanTemplate, "http://") || strings.HasSuffix(envConfig.PlanTemplate, "https://")) {
			resp, err := http.Get(envConfig.PlanTemplate)
			if err != nil {
				log.Fatalf("Failed to read plan template from URL %q : %s\n", envConfig.PlanTemplate, err)
				os.Exit(100)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Failed to read plan template from URL %q : %s\n", envConfig.PlanTemplate, err)
				os.Exit(100)
			}
			planTemplate = string(body)
		} else {
			body, err := ioutil.ReadFile(envConfig.PlanTemplate)
			if err != nil {
				log.Fatalf("Failed to read plan template from file %q : %s\n", envConfig.PlanTemplate, err)
				os.Exit(100)
			}
			planTemplate = string(body)
		}
	}

	// read all go files from the directory that was given. Note that if no directory (GO_PACKAGE)
	// was given, this is ./
	files, err := filepath.Glob(fmt.Sprintf("%s*.go", goPackage))
	if err != nil {
		log.Printf("failed to read directory %s : %s\n", goPackage, err)
		printSupportedFunctionsAndExit(detector.Signatures())
	}

	for _, f := range files {
		log.Printf("Processing file %s\n", f)
		// read file
		file, err := os.Open(f)
		if err != nil {
			log.Println(err)
			os.Exit(100)
		}
		defer file.Close()

		// read the whole file in
		srcbuf, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			printSupportedFunctionsAndExit(detector.Signatures())
		}
		f := &detect.Function{File: f, Source: string(srcbuf)}
		deets, err := detector.CheckFile(f)
		if err != nil {
			log.Panic("Failed to process file %q : %s", f, err)
		}
		if deets != nil {
			log.Printf("Found supported function %q in package %q signature %q", deets.Name, deets.Package, deets.Signature)
			// If the user didn't specify a specific function, use it. If they specified the function, make sure it
			// matches what we found.
			if goFunction == "" || goFunction == deets.Name {
				deets.Package = fullGoPackage
				if err := writePlan(planFileName, planTemplate, deets); err != nil {
					log.Println("failed to write the build plan: ", err)
					os.Exit(100)
				}
				os.Exit(0)
			}
		}
	}
	printSupportedFunctionsAndExit(detector.Signatures())
}

func writePlan(planFileName, planTemplate string, details *detect.FunctionDetails) error {
	planFile, err := os.OpenFile(planFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer planFile.Close()

	t, err := template.New("Plan").Parse(planTemplate)
	if err != nil {
		log.Printf("Failed to parse the template file : %s\n", err)
		return err
	}
	args := PlanArguments{
		Name: defaultPlanName,
		Function: details.Name,
		Package: details.Package,
		Env: make(map [string]string),
	}
	for _, env := range os.Environ() {
		// env pieces are ENV=VALUE, so split them so we get a key=>value into map, which to index.
		pieces := strings.SplitN(env, "=", 2)
		if len(pieces) > 1 {
			args.Env[pieces[0]] = pieces[1]
		}
	}
	if err := t.Execute(planFile, args); err != nil {
		log.Printf("Failed to execute template file: %s\n", err)
		return err
	}
	return nil
}

// readModuleName is a terrible hack for yanking the module from go.mod file.
// Should be replaced with something that actually understands go...
func readModuleName() (string, error) {
	modFile, err := os.Open("./go.mod")
	if err != nil {
		return "", err
	}
	defer modFile.Close()
	scanner := bufio.NewScanner(modFile)
	for scanner.Scan() {
		pieces := strings.Split(scanner.Text(), " ")
		fmt.Printf("FOund pieces as %+v\n", pieces)
		if len(pieces) >= 2 && pieces[0] == "module" {
			return pieces[1], nil
		}
	}
	return "", nil
}
