package client

import "github.com/m-szalik/pact-framework-example/go-api-client/model"

type BookClient interface {
	GetBooks() ([]model.ClientBook, error)
	GetBookByID(id int) (*model.ClientBook, error)
	Delete(id int) error
}
