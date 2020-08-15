package routes

import (
	"github.com/gorilla/mux"
	c "ufc.com/deti/go-dad/src/controller"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/books/{bucket_name}", c.GetAll).Methods("GET")
	router.HandleFunc("/books", c.Store).Methods("POST")
	router.HandleFunc("/books/{bucket_name}/{id}", c.GetOne).Methods("GET")
	router.HandleFunc("/books/{bucket_name}/{id}", c.Delete).Methods("DELETE")
	return router
}
