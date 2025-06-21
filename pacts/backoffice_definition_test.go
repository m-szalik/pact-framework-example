package pacts

import (
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBackofficeDefineBookPact(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "Backoffice",
		Provider: "BooksAPI",
	})
	assert.NoError(t, err)

	t.Run("Create a book", func(t *testing.T) {
		mockProvider.
			AddInteraction().
			Given("create new book").
			WithRequest("POST", "/books", func(builder *consumer.V4RequestBuilder) {
				builder.JSONBody(Book{ID: 11, Title: "New book"})
			}).
			WillRespondWith(201, func(builder *consumer.V4ResponseBuilder) {
				builder.JSONBody(Book{ID: 11, Title: "New book"})
			})
		err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			resp, err := request(t, config, "POST", "/books", Book{ID: 11, Title: "New book"}, func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			})
			assert.Equal(t, 201, resp.StatusCode)
			assert.NoError(t, err)
			return nil
		})
		assert.NoError(t, err)
	})

}
