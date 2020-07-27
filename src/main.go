package main

import (
	"net/http"

	r "ufc.com/deti/go-dad/src/routes"
)

func main() {
	route := r.NewRouter()
	http.ListenAndServe(":8080", route)
}
