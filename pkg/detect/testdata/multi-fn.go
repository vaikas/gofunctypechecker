package function

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"net/http"
)

func ReceiveHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("doing stuff here")
}

func ReceiveEvent(ce cloudevents.Event) (*cloudevents.Event, error) {
	r := cloudevents.NewEvent(cloudevents.VersionV1)
	r.SetType("io.mattmoor.cloudevents-go-fn")
	r.SetSource("https://github.com/mattmoor/cloudevents-go-fn")

	if err := r.SetData("application/json", struct {
		A string `json:"a"`
		B string `json:"b"`
	}{
		A: "hello",
		B: "world",
	}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}

	return &r, nil
}
