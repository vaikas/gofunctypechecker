package detect

import "testing"

const signatureFileJSON = "./testdata/signatures.json"
const signatureFileYAML = "./testdata/signatures.yaml"
const signatureFileTOML = "./testdata/signatures.toml"

var tests = map[string]*FunctionDetails{
	"./testdata/f1.go": &FunctionDetails{
		Name:      "Receive",
		Signature: "func(http.ResponseWriter, *http.Request)",
	},
	"./testdata/f2.go": &FunctionDetails{
		Name:      "Receive2",
		Signature: "func(http.ResponseWriter, *http.Request)",
	},
	"./testdata/f3-bad.go": nil,
	"./testdata/f4-bad.go": nil,
	"./testdata/f5.go": &FunctionDetails{
		Name:      "Receive4",
		Signature: "func(context.Context, v2.Event) (*v2.Event, error)",
	},
	"./testdata/f6.go": &FunctionDetails{
		Name:      "Receive5",
		Signature: "func(v2.Event) (*v2.Event, error)",
	},
	"./testdata/f7.go": &FunctionDetails{
		Name:      "Impl",
		Signature: "func(context.Context, <-chan *proto.Request, chan *proto.Response) error",
	},
}

func TestAllCases(t *testing.T) {
	// Valid function signatures
	// func(http.ResponseWriter, *http.Request)
	// func(v2.Event) (*v2.Event, error)
	var validFunctions = []FunctionSignature{
		{In: []FunctionArg{
			{ImportPath: "net/http", Name: "ResponseWriter"},
			{ImportPath: "net/http", Name: "Request", Pointer: true},
		}, Out: []FunctionArg{}},
		{In: []FunctionArg{
			{ImportPath: "context", Name: "Context"},
			{ImportPath: "github.com/cloudevents/sdk-go/v2", Name: "Event"},
		}, Out: []FunctionArg{
			{ImportPath: "github.com/cloudevents/sdk-go/v2", Name: "Event", Pointer: true},
			{Name: "error"},
		}},
		{In: []FunctionArg{
			{ImportPath: "github.com/cloudevents/sdk-go/v2", Name: "Event"},
		}, Out: []FunctionArg{
			{ImportPath: "github.com/cloudevents/sdk-go/v2", Name: "Event", Pointer: true},
			{Name: "error"},
		}},
		{In: []FunctionArg{
			{ImportPath: "context", Name: "Context"},
			{ImportPath: "github.com/mattmoor/korpc-sample/gen/proto", Name: "Request", Pointer: true, Channel: Receive},
			{ImportPath: "github.com/mattmoor/korpc-sample/gen/proto", Name: "Response", Pointer: true, Channel: Both},
		}, Out: []FunctionArg{
			{Name: "error"},
		}},
	}

	d := NewDetector(validFunctions)
	runTests(t, d)
}

func TestAllCasesFromJSONFile(t *testing.T) {
	d, err := NewDetectorFromFile(signatureFileJSON)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileJSON, err)
	}
	runTests(t, d)
}

func TestAllCasesFromYAMLFile(t *testing.T) {
	d, err := NewDetectorFromFile(signatureFileYAML)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileJSON, err)
	}
	runTests(t, d)
}

func TestAllCasesFromTOMLFile(t *testing.T) {
	d, err := NewDetectorFromFile(signatureFileTOML)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileJSON, err)
	}
	runTests(t, d)
}

func runTests(t *testing.T, d *Detector) {
	for file, expected := range tests {
		t.Run(file, func(t *testing.T) {
			got, err := d.ReadAndCheckFile(file)
			if err != nil {
				t.Fatalf("Failed to check file: %s", err)
			}
			if expected != nil && got == nil {
				t.Errorf("Expected %+v but got nil", expected)
			} else if expected == nil && got != nil {
				t.Errorf("Expected nil but got %+v", got)
			} else if expected != nil && got != nil {
				if expected.Name != got.Name {
					t.Errorf("Function Name differs got %q expected %q", got.Name, expected.Name)
				}
				if expected.Signature != got.Signature {
					t.Errorf("Signature differs got %q expected %q", got.Signature, expected.Signature)
				}
			}
		})
	}
}
