package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	r "ufc.com/deti/go-dad/src/routes"	
)

func main() {
	mux := r.NewRouter()
	c := cors.New(cors.Options{
		AllowedMethods: []string{"POST", "GET", "DELETE", "PUT", "PATCH"},
	})
	handler := c.Handler(mux)
	fmt.Println("Server running on 80 port ... ")
	log.Fatal(http.ListenAndServe(":80", handler))
}
