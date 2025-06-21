package main

import (
	"fmt"
	"github.com/m-szalik/pact-framework-example/go-api-client/client"
)

func main() {
	bookClient := client.NewHTTPBookClient("http://localhost:8080")
	books, err := bookClient.GetBooks()
	if err != nil {
		panic(fmt.Sprintf("Error fetching books: %s\n", err))
	}
	fmt.Println("All Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s\n", book.ID, book.Title)
	}
}
