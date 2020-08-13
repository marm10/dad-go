package model

import (
	"errors"
)

type Book struct {
	Name    string   `json:"nome"`
	Authors []string `json:"autores"`
	Year    string   `json:"data_lancamento"`
	Preco   float64  `json:"preco"`
	Id      int      `json:"id"`
	Cover   string   `json:"nome_capa"`
}

type Books []Book

var books = Books{}

func GetAll() Books {
	return books
}

func GetOne(id int) (*Book, error) {

	err := errors.New("Book not found")
	for _, book := range books {
		if book.Id == id {
			return &book, nil
		}
	}
	return nil, err

}

func Store(book Book) *Book {
	book.Id = len(books) + 1
	books = append(books, book)
	return &book
}

func Delete(id int) error {

	err := errors.New("Book not found")
	for index, book := range books {
		if book.Id == id {
			books = append(books[:index], books[index+1:]...)
			return nil
		}
	}
	return err
}
