#  PactFramework example - Go API Client

## Discretionary 

 * Consumer - for example webapp
 * Provider - backend


## Problem
To test the client we can run the server and try to connect to the server, but what versions of the server to run?
And first of all running that server is not so easy as likely there is database and other components.


## Solution
You can start your own mock server and test your client against it.
But that's your understanding of server API.
[Example](client/http_client_test.go)

But what about versioning of your API?

## Solution 2 - PactFramework
```shell
go get github.com/pact-foundation/pact-go/v2
go install github.com/pact-foundation/pact-go/v2
pact-go -l DEBUG install
```


## Push pacts to broker
```shell
pact-broker publish --broker-base-url=http://localhost:9292 --consumer-app-version=2 ./client/pacts
```