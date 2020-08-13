package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	h "ufc.com/deti/go-dad/src/handlerException"
	b "ufc.com/deti/go-dad/src/model"	

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	BUCKET_NAME = "book-covers"
	REGION = "us-east-2"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	books := b.GetAll()
	if err := json.NewEncoder(w).Encode(&books); err != nil {
		h.Handler(w, r, http.StatusInternalServerError, err.Error())
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	book, err := b.GetOne(id)
	if err != nil {
		h.Handler(w, r, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&book)

}

func Store(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var book b.Book
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&book); err != nil {
		h.Handler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	bookCreated := b.Store(book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookCreated)

}

func Delete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	err := b.Delete(id)
	if err != nil {
		h.Handler(w, r, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func CreateBucket(w http.ResponseWriter, r *http.Request) () {
	w.Header().Set("Content-Type", "application/json")

	
	// snippet-start:[s3.go.create_bucket.call]
    svc := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))

    _, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(BUCKET_NAME),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION)
		},
    })
    // snippet-end:[s3.go.create_bucket.call]

    // snippet-start:[s3.go.create_bucket.wait]
	
	if err != nil {
		h.Handler(w, r, http.StatusBadRequest, err.Error())
		return
	}
    // snippet-end:[s3.go.create_bucket.wait]

	w.WriteHeader(http.StatusAccepted)	
}
