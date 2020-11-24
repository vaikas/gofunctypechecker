package detect

import "testing"

const signatureFile = "./testdata/signatures.json"

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
	}
	tests := map[string]*FunctionDetails{
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
		"./testdata/f5.go":     &FunctionDetails{
			Name: "Receive4",
			Signature: "func(context.Context, v2.Event) (*v2.Event, error)",
		},
		"./testdata/f6.go": &FunctionDetails{
			Name:      "Receive5",
			Signature: "func(v2.Event) (*v2.Event, error)",
		},
	}

	d := NewDetector(validFunctions)
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

func TestAllCasesFromFile(t *testing.T) {
	tests := map[string]*FunctionDetails{
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
		"./testdata/f5.go":     &FunctionDetails{
			Name: "Receive4",
			Signature: "func(context.Context, v2.Event) (*v2.Event, error)",
		},
		"./testdata/f6.go": &FunctionDetails{
			Name:      "Receive5",
			Signature: "func(v2.Event) (*v2.Event, error)",
		},
	}

	d, err := NewDetectorFromFile(signatureFile)
	if err != nil {
		t.Fatalf("Failed to read function signatures from file: %q : %s ", signatureFile, err)
	}
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
