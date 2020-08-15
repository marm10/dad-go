package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"strings"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"
	h "ufc.com/deti/go-dad/src/handlerException"
	b "ufc.com/deti/go-dad/src/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	BUCKET_NAME = "book-covers-dad"
	REGION = "us-east-2"
)

var (
	s3session *s3.S3
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	att := mux.Vars(r)
	bucketName := att["bucket_name"]

	fmt.Println("bucketname")
	fmt.Println(bucketName)

	result, _ := listObjects(bucketName)
	contents := result.Contents

	var books []b.Book

	for i, s := range contents {
		fmt.Println(i)
		var book b.Book

		if strings.Contains(aws.StringValue(s.Key), ".json") {
			object, _ := GetObject(aws.StringValue(s.Key), bucketName)
			readarr := bytes.NewReader(object)
			decoder := json.NewDecoder(readarr)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&book); err != nil {
				h.Handler(w, r, http.StatusBadRequest, err.Error())
				return
			}

			books = append(books, book)
		}
	}

	if err := json.NewEncoder(w).Encode(&books); err != nil {
		h.Handler(w, r, http.StatusInternalServerError, err.Error())
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	att := mux.Vars(r)
	bucketName := att["bucket_name"]
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	
	returnBook, err := GetBookById(id, bucketName)

	if err != nil {
		h.Handler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&returnBook)

}

func GetBookById(id int, bucketName string) ( returnBook b.Book, err error)  {
	fmt.Println("bucketname")
	fmt.Println(bucketName)

	result, _ := listObjects(bucketName)
	contents := result.Contents

	var books []b.Book
	var book b.Book

	for i, s := range contents {
		fmt.Println(i)

		if strings.Contains(aws.StringValue(s.Key), ".json") {
			object, _ := GetObject(aws.StringValue(s.Key), bucketName)
			readarr := bytes.NewReader(object)
			decoder := json.NewDecoder(readarr)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&book)

			if book.Id == id {
				returnBook = book
			}
		}
	}

	return returnBook, err
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

	CreateBucket(bucket_name)

	price, _ := strconv.ParseFloat(priceAttr, 64)
	authors := strings.Split(authorsAttr, ",")

	book.Name = name
	book.Authors = authors
	book.Year = year
	book.Preco = price
	book.Cover = fileName;

	objects, err := listObjects(bucket_name)

	if err != nil {
		h.Handler(w, r, http.StatusNotFound, err.Error())
		return
	}

	contents := result.Contents
	book.Id = len(contents)

	folderName := "book_"+name+"/";

	_, err1 := s3session.PutObject(&s3.PutObjectInput{
		Body: file,
		Bucket: aws.String(bucket_name),
		Key: aws.String(folderName+fileName),
		ACL: aws.String(s3.BucketCannedACLPublicRead),
	})
	
	if err1 != nil {
		h.Handler(w, r, http.StatusNotFound, err1.Error())
		return
	}
	
	defer f.Close()

	b, _ := json.Marshal(&book)

	br := bytes.NewReader(b)

	w.WriteHeader(http.StatusCreated)
	_, err2 := s3session.PutObject(&s3.PutObjectInput{
		Body: br,
		Bucket: aws.String(bucket_name),
		Key: aws.String(folderName+"book.json"),
		ACL: aws.String(s3.BucketCannedACLPublicRead),
	})

	if err2 != nil {
		panic(err2)
	}

	json.NewEncoder(w).Encode(&book)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	att := mux.Vars(r)
	idAtt := att["id"]
	id, _ := strconv.Atoi(idAtt)
	bucketName := att["bucket_name"]

	returnBook, err := GetBookById(id, bucketName)

	if err != nil {
		h.Handler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := DeleteObject(returnBook.Name)
	
	if err != nil {
		h.Handler(w, r, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func init() {
	if s3session == nil {
		s3session = s3.New(session.Must(session.NewSession(&aws.Config{
			Region: aws.String(REGION),
		})))
	}	
}

func CreateBucket(bucket_name string) (err error) {
    _, err0 := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket_name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
	})
	
	if err0 != nil {
		if aerr, ok := err0.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already exists!")
				err = err0
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket name exists and is owned by you!")
			default:
				err = err0
			}
		}
	}
}

func listObjects(bucketName string) (resp *s3.ListObjectsV2Output, err error) {
	return s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
} 

func GetObject(fileName string, bucketName string) (obj []byte, err error){
	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
	})

	obj, err1 := ioutil.ReadAll(resp.Body)

	if err1 != nil {
		panic(err1)
	}
}

func DeleteObject(bucketName string, fileName string) (resp *s3.DeleteObjectOutput, err error) {
	return s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
	})
}
