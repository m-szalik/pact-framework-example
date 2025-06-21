package pacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"io"
	"net/http"
	"testing"
)

func request(t *testing.T, config consumer.MockServerConfig, method, url string, body any, mods ...func(r *http.Request)) (*http.Response, error) {
	var reqBody io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reqBody = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d%s", config.Host, config.Port, url), reqBody)
	if err != nil {
		return nil, err
	}
	for _, mod := range mods {
		mod(req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
