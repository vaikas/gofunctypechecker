package function

import (
	ctx "context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/protocol"
	cloudevents "github.com/someothercloudevents/sdk-go/v2"
)

func Receive3(ctx ctx.Context, event cloudevents.Event) (*cloudevents.Event, protocol.Result) {
	fmt.Printf("☁️  CloudEvents.Event\n%s", event.String())
	return nil, nil
}
