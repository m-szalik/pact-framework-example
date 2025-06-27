package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/m-szalik/goutils"
	"github.com/m-szalik/pact-framework-example/go-api-client/model"
)

// hTTPBookClient implements BookClient using HTTP
type hTTPBookClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *hTTPBookClient) Delete(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/books/%d", c.BaseURL, id), nil)
	if err != nil {
		return err
	}
	_, err = c.HTTPClient.Do(req)
	return err
}

// NewHTTPBookClient returns an initialized hTTPBookClient
func NewHTTPBookClient(baseURL string) BookClient {
	return &hTTPBookClient{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *hTTPBookClient) GetBooks() ([]model.ClientBook, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/books")
	if err != nil {
		return nil, err
	}
	defer goutils.CloseQuietly(resp.Body)

	var books []model.ClientBook
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		return nil, err
	}
	return books, nil
}

func (c *hTTPBookClient) GetBookByID(id int) (*model.ClientBook, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/books/%d", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer goutils.CloseQuietly(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("book with ID %d not found", id)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var book model.ClientBook
	if err := json.Unmarshal(body, &book); err != nil {
		return nil, err
	}
	return &book, nil
}
