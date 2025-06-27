package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
)

func TestVerifyServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_server := newServer()
	go _server.serverStart(ctx)

	verifier := provider.NewVerifier()
	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("Authorization", "Bearer testingToken")
			next.ServeHTTP(w, r)
		})
	}

	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL:    "http://127.0.0.1:8080",
		Provider:           "BooksAPIProvider",
		FailIfNoPactsFound: true,
		ProviderVersion:    "latestVersion",
		BrokerURL:          "http://localhost:9292",
		// ConsumerVersionSelectors: []provider.Selector{
		//	&provider.ConsumerVersionSelector{
		//		Tag: "develop",
		//	},
		//	&provider.ConsumerVersionSelector{
		//		Tag: "develop",
		//	},
		// },
		PublishVerificationResults: true,
		RequestFilter:              f,
	})
	assert.NoError(t, err)

}
