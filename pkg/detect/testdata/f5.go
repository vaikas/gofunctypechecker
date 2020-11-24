package function

import (
	ctx "context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func Receive4(ctx ctx.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	fmt.Printf("☁️  CloudEvents.Event\n%s", event.String())
	return nil, nil
}
