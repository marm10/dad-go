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
	BUCKET_NAME = "book-covers-dad"
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

	file, handler, err := r.FormFile("file")
	fileName := r.FormValue("file_name")
    if err != nil {
        panic(err)
	}
	
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        panic(err)
	}
	
	defer f.Close()

	w.Header().Set("Content-Type", "application/json")
	var book b.Book
	book.Name = r.FormValue("nome")
	book.Authors = r.FormValue("autores")
	book.Year = r.FormValue("data_lancamento")
	book.Preco = r.FormValue("preco")

	UploadObject(file, filename)

	bookCreated := b.Store(book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookCreated)

}

func UploadObject(File f, string filename) {
	svc := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))

	_, err := svc.PutObject(&s3.PutObjectInput{
		Body: f,
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(filename),
		ACL: aws.String(s3.BucketCannedACLPublicRead)
	})
	
	if err != nil {
		panic(err)
	}

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
			LocationConstraint: aws.String(REGION),
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
