package pacts

import (
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebappDefineBookPact(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "Webapp",
		Provider: "BooksAPI",
	})
	assert.NoError(t, err)

	t.Run("Get all books", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("list of books is not empty").
			UponReceiving("A request for list of books").
			WithRequest("GET", "/books").
			WillRespondWith(200, func(builder *consumer.V4ResponseBuilder) {
				builder.Header("Content-Type", matchers.S("application/json"))
				builder.JSONBody([]Book{{ID: 5, Title: "Effective Java"}, {ID: 6, Title: "Refactoring"}})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(t, config, "GET", "/books", nil)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Get exiting book by ID", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 5 exists").
			UponReceiving("A request for Book 5").
			WithRequest("GET", "/books/5").
			WillRespondWith(200, func(builder *consumer.V4ResponseBuilder) {
				builder.Header("Content-Type", matchers.S("application/json"))
				builder.JSONBody(Book{ID: 5, Title: "Effective Java"})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(t, config, "GET", "/books/5", nil)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Get not exiting book by ID", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 99 does not exist").
			UponReceiving("A request for Book 99").
			WithRequest("GET", "/books/99").
			WillRespondWith(404, func(builder *consumer.V4ResponseBuilder) {
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(t, config, "GET", "/books/99", nil)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Delete exiting book", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 5 exists").
			UponReceiving("Delete request for Book 5").
			WithRequest("DELETE", "/books/5", func(builder *consumer.V4RequestBuilder) {
				// additional request conditions here
			}).
			WillRespondWith(200, func(builder *consumer.V4ResponseBuilder) {
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			_, err := request(t, config, http.MethodDelete, "/books/5", nil)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

}
