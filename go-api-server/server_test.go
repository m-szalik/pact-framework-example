package main

import (
	"encoding/json"
	"fmt"
	"github.com/m-szalik/pact-framework-example/go-api-server/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerHttp(t *testing.T) {
	s := newServer()
	handler := s.httpHandler()
	t.Run("getList of books", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/books", nil)
		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
		buff, err := io.ReadAll(writer.Result().Body)
		assert.NoError(t, err)
		var data []model.ServerBook
		err = json.Unmarshal(buff, &data)
		assert.NoError(t, err)
		for _, book := range data {
			fmt.Printf(" * ServerBook: %+v\n", book)
		}
		// TODO more asserts here....
	})

	t.Run("book by id - not found", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/book/999", nil)
		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)
	})

}
