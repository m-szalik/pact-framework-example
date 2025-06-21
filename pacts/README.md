#  PactFramework example - Contract definitions

# Pact Broker

```shell
docker run --rm -d -e PACT_BROKER_DATABASE_ADAPTER="sqlite" -e PACT_BROKER_DATABASE_URL="sqlite:////tmp/pact_broker.sqlte3" -p 9292:9292 pactfoundation/pact-broker:latest
```
http://localhost:9292



## Push pacts to broker
```shell
pact-broker publish --broker-base-url=http://localhost:9292 --consumer-app-version=2 ./pacts
```

```shell
export VER=`git rev-parse --short HEAD`
pact-broker publish ./pacts/ --broker-base-url=http://localhost:9292 --consumer-app-version=${VER}  -t "develop"
```
