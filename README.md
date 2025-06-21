# Contract testing:

## List of other frameworks 
   1. Pact - widely adopted, multiple language support.
   2. OpenAPI (Swagger) 
   3. Dredd - OpenAPI contract validation

## What is contract

## Who define contract
 * Consumer-driven
 * Producer-driven
 * Schema/Message Based

# Why Pact
1. Responsive community.
2. Multiple language support
3. Easy to use
4. Tools
5. Unfortunately it is consumer-driven


## Example
1. Go Consumer (Client)
2. Go provider (Server)
3. Pacts and Pact broker

# Benefits of contract testing
1. You can organize your code that all contract are in a separated repo, so spacial approval policy can be applied.
2. Any contract violation (on the provider or consumer) will be discovered impliedly.
3. You need to start design from the contract --> API first approach
4. API Mocks out of the box











# CLI Tools
https://docs.pact.io/implementation_guides/cli#pact-broker-cli

```shell
curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | PACT_CLI_VERSION=v2.4.25 bash
```