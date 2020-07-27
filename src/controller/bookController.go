package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	h "ufc.com/deti/go-dad/src/handlerException"
	b "ufc.com/deti/go-dad/src/model"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	books := b.GetAll()
	if err := json.NewEncoder(w).Encode(books); err != nil {
		panic(err)
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	book, err := b.GetOne(id)
	if err != nil {
		errorDefault := h.DefaultError{
			StatusCode: 404,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(&errorDefault)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book b.Book
	json.NewDecoder(r.Body).Decode(&book)
	err := b.Store(book)
	if err != nil {
		json.NewEncoder(w).Encode(&book)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	err := b.Delete(id)
	if err != nil {
		errorDefault := h.DefaultError{
			StatusCode: 404,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(&errorDefault)
	}
}
