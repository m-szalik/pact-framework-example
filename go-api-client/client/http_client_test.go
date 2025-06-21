package client

import (
	_ "embed"
	"fmt"
	"github.com/m-szalik/pact-framework-example/go-api-client/model"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
)

//go:embed books-response.json
var booksResponseJson string

// TestWithMockHttpServer - test using httptest without contract testing
func TestWithMockHttpServer(t *testing.T) {
	smux := http.NewServeMux()
	smux.HandleFunc("/books", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(booksResponseJson))
	})
	svr := httptest.NewServer(smux)
	defer svr.Close()
	slog.Info(fmt.Sprintf("Test server started under %s", svr.URL))

	bookClient := NewHTTPBookClient(svr.URL)
	books, err := bookClient.GetBooks()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(books))
	fmt.Printf("Books: %+v\n", books)
}

// TestWithPact test with pact, that generates pact contract
func TestWithPact(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "BooksAPIConsumer",
		Provider: "BooksAPIProvider",
	})
	assert.NoError(t, err)
	mockProvider.
		AddInteraction().
		Given("A book with ID 5 exists").
		UponReceiving("A request for Book 5").
		WithRequest("GET", "/books/5").
		WillRespondWith(200, func(builder *consumer.V4ResponseBuilder) {
			builder.Header("Content-Type", matchers.S("application/json"))
			builder.JSONBody(model.ClientBook{ID: 5, Title: "Effective Java"})
		})
	err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		bookClient := NewHTTPBookClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
		book, err := bookClient.GetBookByID(5)
		assert.NoError(t, err)
		assert.Equal(t, 5, book.ID)
		assert.Equal(t, "Effective Java", book.Title)
		return nil
	})
	assert.NoError(t, err)

	mockProvider.
		AddInteraction().
		Given("A book with ID 5 does not exist").
		UponReceiving("A request for Book 5").
		WithRequest("GET", "/books/5").
		WillRespondWith(404, func(builder *consumer.V4ResponseBuilder) {})
	err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		bookClient := NewHTTPBookClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
		book, err := bookClient.GetBookByID(5)
		assert.Error(t, err)
		assert.Nil(t, book)
		return nil
	})
	assert.NoError(t, err)
}

// TestWithPactNewEndpointDelete - another test with pact that will add new definition to the contract.
func TestWithPactNewEndpointDelete(t *testing.T) {
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "BooksAPIConsumer",
		Provider: "BooksAPIProvider",
	})
	assert.NoError(t, err)

	t.Run("Delete book that exits", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 5 exists").
			UponReceiving("Delete request for Book 5").
			WithRequest("DELETE", "/books/5").
			WillRespondWith(200, func(builder *consumer.V2ResponseBuilder) {})
		err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			bookClient := NewHTTPBookClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			err := bookClient.Delete(5)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Delete book that does not exist", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 5 does not exist").
			UponReceiving("A request for Book 5").
			WithRequest("GET", "/books/5").
			WillRespondWith(404, func(builder *consumer.V2ResponseBuilder) {})
		err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			bookClient := NewHTTPBookClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			book, err := bookClient.GetBookByID(5)
			assert.Error(t, err)
			assert.Nil(t, book)
			return nil
		})
		assert.NoError(t, err)
	})

}
