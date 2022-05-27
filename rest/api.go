package rest

import (
	"fmt"
	"net/http"
)

func HandleHealth(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}
