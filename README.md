# Contract testing:

## List of other frameworks 
   1. Pact - widely adopted, multiple language support.
   2. OpenAPI (Swagger) 
   3. Dredd - OpenAPI contract validation

## What is contract

## Who defines contract
 * Consumer-driven
 * Producer-driven
 * Schema/Message based

# Why Pact
1. Responsive community.
2. Multiple language support.
3. Easy to use.
4. Supports Rest, grpc and messages.
5. Tools.
6. Unfortunately, it is consumer-driven.


## Example
1. Go Consumer (Client)
2. Go provider (Server)
3. Pacts and Pact broker

## Tools:
### Broker

### Pact-stub-service
```shell
pact-stub-service --port=8080 pacts/*.json
```
```shell
curl -sv http://localhost:8080/books |jq
```
Alternative: https://github.com/pact-foundation/pact-stub-server
or
```shell
docker run --rm -it -v "$(pwd)/pacts:/pacts" -p 8080:8080 pactfoundation/pact-stub-server --dir /pacts --port=8080
```

### pact-provider-verifier


# Benefits of contract testing
1. You can organize your code that all contracts are in a separated repo, so spacial approval policy can be applied.
2. Any contract violation (on the provider or consumer) will be discovered impliedly.
3. You need to start design from the contract → API first approach
4. API Mocks out of the box











# CLI Tools
https://docs.pact.io/implementation_guides/cli#pact-broker-cli

```shell
curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | PACT_CLI_VERSION=v2.4.25 bash
```