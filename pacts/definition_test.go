package pacts

import (
	"bytes"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDefineBookPact(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "BooksAPIConsumer",
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
			resp, err := request(config, "GET", "/books")
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
				builder.JSONBody(Book{ID: 5, Title: "Cool title"})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(config, "GET", "/books/5")
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Get not exiting book by ID", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("A book with ID 5 does not exist").
			UponReceiving("A request for Book 5").
			WithRequest("GET", "/books/5").
			WillRespondWith(404, func(builder *consumer.V4ResponseBuilder) {
				builder.BinaryBody([]byte{1})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(config, "GET", "/books/5")
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
				builder.Query("q", matchers.S("qqq"))
			}).
			WillRespondWith(200, func(builder *consumer.V4ResponseBuilder) {
				builder.BinaryBody([]byte{112})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			_, err := request(config, http.MethodDelete, "/books/5?q=qqq")
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

}

func request(config consumer.MockServerConfig, method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d%s", config.Host, config.Port, url), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
