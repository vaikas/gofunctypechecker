package function

import (
	"fmt"

	ce "github.com/cloudevents/sdk-go/v2"
)

func Receive(event ce.Event) {
	fmt.Printf("☁️  CloudEvents.Event\n%s", event.String())
}
