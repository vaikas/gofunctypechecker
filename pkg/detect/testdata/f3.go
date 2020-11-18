package function

import (
	ctx "context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

func Receive3(ctx ctx.Context, event cloudevents.Event) (*cloudevents.Event, protocol.Result) {
	fmt.Printf("☁️  CloudEvents.Event\n%s", event.String())
	return nil, nil
}
