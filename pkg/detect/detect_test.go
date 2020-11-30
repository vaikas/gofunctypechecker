package detect

import "testing"

const signatureFileJSON = "./testdata/signatures.json"
const signatureFileYAML = "./testdata/signatures.yaml"
const signatureFileTOML = "./testdata/signatures.toml"
const signatureFileURLJSON = "https://raw.githubusercontent.com/vaikas/gofunctypechecker/main/pkg/detect/testdata/signatures.json"
const signatureFileURLYAML = "https://raw.githubusercontent.com/vaikas/gofunctypechecker/main/pkg/detect/testdata/signatures.yaml"
const signatureFileURLTOML = "https://raw.githubusercontent.com/vaikas/gofunctypechecker/main/pkg/detect/testdata/signatures.toml"

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
	"./testdata/f8.go": &FunctionDetails{
		Name:      "Receive",
		Signature: "func(http.ResponseWriter, *http.Request)",
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
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileYAML, err)
	}
	runTests(t, d)
}

func TestAllCasesFromTOMLFile(t *testing.T) {
	d, err := NewDetectorFromFile(signatureFileTOML)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileTOML, err)
	}
	runTests(t, d)
}

func TestAllCasesFromURLJSON(t *testing.T) {
	d, err := NewDetectorFromURL(signatureFileURLJSON)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileURLJSON, err)
	}
	runTests(t, d)
}

func TestAllCasesFromURLYAML(t *testing.T) {
	d, err := NewDetectorFromURL(signatureFileURLYAML)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileURLYAML, err)
	}
	runTests(t, d)
}

func TestAllCasesFromTOML(t *testing.T) {
	d, err := NewDetectorFromURL(signatureFileURLTOML)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileURLTOML, err)
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

func TestMultiDetect(t *testing.T) {
	d, err := NewDetectorFromFile(signatureFileJSON)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFileJSON, err)
	}

	want := []FunctionDetails{
		{
			Name:      "ReceiveHTTP",
			Signature: "func(http.ResponseWriter, *http.Request)",
		},
		{
			Name:      "ReceiveEvent",
			Signature: "func(v2.Event) (*v2.Event, error)",
		},
	}
	testfile := "./testdata/multi-fn.go"

	if got, err := d.ReadAndCheckFile(testfile); err == nil {
		t.Errorf("Expected checking error on multi-file, got: %+v", got)
	}

	got, err := d.ReadAllFromFile(testfile)
	if err != nil {
		t.Errorf("Expected no error detecting signatures, got %v", err)
	}
	if len(got) != len(want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("Error at %d, wanted %v, got %v", i, want[i], got[i])
		}
	}
}
