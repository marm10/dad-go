package routes

import (
	"github.com/gorilla/mux"
	c "ufc.com/deti/go-dad/src/controller"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/books", c.GetAll).Methods("GET")
	router.HandleFunc("/createBucket", c.CreateBucket).Methods("GET")
	router.HandleFunc("/books", c.Store).Methods("POST")
	router.HandleFunc("/books/{id}", c.GetOne).Methods("GET")
	router.HandleFunc("/books/{id}", c.Delete).Methods("DELETE")
	return router
}
