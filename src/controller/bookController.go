package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"strings"
	"bytes"
	"fmt"

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

	file, handler, err := r.FormFile("capa")
	fileName := r.FormValue("nome_capa")
    if err != nil {
        panic(err)
	}
	
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        panic(err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	var book b.Book
	name := r.FormValue("nome")
	authorsAttr := r.FormValue("autores")
	year := r.FormValue("data_lancamento")
	priceAttr := r.FormValue("preco")
	bucket_name := r.FormValue("nome_bucket")

	s3session := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))

	CreateBucket(bucket_name, s3session)

	price, _ := strconv.ParseFloat(priceAttr, 64)
	authors := strings.Split(authorsAttr, ",")

	book.Name = name
	book.Authors = authors
	book.Year = year
	book.Preco = price
	book.Cover = fileName;

	bookCreated := b.Store(book)

	folderName := "book_"+name+"/";

	_, err1 := s3session.PutObject(&s3.PutObjectInput{
		Body: file,
		Bucket: aws.String(bucket_name),
		Key: aws.String(folderName+fileName),
		ACL: aws.String(s3.BucketCannedACLPublicRead),
	})
	
	if err1 != nil {
		panic(err1)
	}
	
	defer f.Close()

	b, _ := json.Marshal(bookCreated)

	br := bytes.NewReader(b)

	w.WriteHeader(http.StatusCreated)
	_, err2 := ss3sessionvc.PutObject(&s3.PutObjectInput{
		Body: br,
		Bucket: aws.String(bucket_name),
		Key: aws.String(folderName+"book.json"),
		ACL: aws.String(s3.BucketCannedACLPublicRead),
	})

	if err2 != nil {
		panic(err2)
	}

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

func CreateBucket(bucket_name string, s3session *s3.S3) () {
    _, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket_name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
    })
	
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already exists!")
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket name exists and is owned by you!")
			default:
				panic(err)	
			}
		}
	}
    // snippet-end:[s3.go.create_bucket.wait]	
}
