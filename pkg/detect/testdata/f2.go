package functionpkg

import (
	"fmt"
	nethttp "net/http"
)

func Receive2(writer nethttp.ResponseWriter, request *nethttp.Request) {
	fmt.Println("doing stuff here")
}
