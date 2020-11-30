package function

import (
	"fmt"
	"net/http"
)

func Receive(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("doing stuff here")
}

func HaHa(s string) {
	fmt.Println("Nelson says 'HA HA!' at " + s)
}
