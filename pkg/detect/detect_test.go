package detect

import "testing"

func TestAllCases(t *testing.T) {

	tests := map[string]*FunctionDetails{
		"./testdata/f1.go": &FunctionDetails{
			Package:   "function",
			Signature: "func(event.Event)",
		},
		"./testdata/f2.go": &FunctionDetails{
			Package:   "functionpkg",
			Signature: "func(context.Context, event.Event) (*event.Event, protocol.Result)",
		},
		"./testdata/f3.go": &FunctionDetails{
			Package:   "function",
			Signature: "func(context.Context, event.Event) (*event.Event, protocol.Result)",
		},
		"./testdata/f4.go": nil,
	}

	for file, expected := range tests {
		t.Run(file, func(t *testing.T) {
			got := ReadAndCheckFile(file)
			if expected != nil && got == nil {
				t.Errorf("Expected %+v but got nil", expected)
			} else if expected == nil && got != nil {
				t.Errorf("Expected nil but got %+v", got)
			} else if expected != nil && got != nil {
				if expected.Package != got.Package {
					t.Errorf("Package differs got %q expected %q", got.Package, expected.Package)
				}
				if expected.Signature != got.Signature {
					t.Errorf("Signature differs got %q expected %q", got.Signature, expected.Signature)
				}
			}
		})
	}
}
