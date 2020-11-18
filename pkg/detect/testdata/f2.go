package functionpkg

import (
	"context"
	"fmt"

	v2 "github.com/cloudevents/sdk-go/v2"
	prot "github.com/cloudevents/sdk-go/v2/protocol"
)

func Receive2(ctx context.Context, event v2.Event) (*v2.Event, prot.Result) {
	fmt.Printf("☁️  CloudEvents.Event\n%s", event.String())
	return nil, nil
}
