package detect

import "testing"

func TestAllCases(t *testing.T) {

	tests := map[string]*FunctionDetails{
		"./testdata/f1.go": &FunctionDetails{
			Name: "Receive",
			Signature: "func(http.ResponseWriter, *http.Request)",
		},
		"./testdata/f2.go": &FunctionDetails{
			Name: "Receive2",
			Signature: "func(http.ResponseWriter, *http.Request)",
		},
		"./testdata/f3-bad.go": nil,
		"./testdata/f4-bad.go": nil,
	}

	for file, expected := range tests {
		t.Run(file, func(t *testing.T) {
			got := ReadAndCheckFile(file)
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
