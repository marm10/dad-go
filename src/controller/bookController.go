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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	books := b.GetAll()
	if err := json.NewEncoder(w).Encode(&books); err != nil {
		errorDefault := h.DefaultError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(&errorDefault)
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	book, err := b.GetOne(id)
	if err != nil {
		errorDefault := h.DefaultError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&errorDefault)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&book)
	}
}

func Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book b.Book
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&book); err != nil {
		errorDefault := h.DefaultError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&errorDefault)
		return
	}
	b.Store(book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&book)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	err := b.Delete(id)
	if err != nil {
		errorDefault := h.DefaultError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&errorDefault)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
