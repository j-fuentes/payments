# Payments API

## How to run

*Dev requirements*:  (code generation from specs), Go >= 1.11.

You can run it natively (requires Go >= 1.11):

```
make serve
```

Or using docker:

```
docker-compose up
```

## How to build

Just run `make` and the binary will be generated in `./out/payments-service`

## Swagger specs

The specs are defined in [./api/swagger.yml](./api/swagger.yml) with Swagger (aka OpenAPI).

Some code is generated from that by running `make generate`. An extra dependency is needed if you want to do so: [go-swagger](https://github.com/go-swagger/go-swagger).

Among the generated code you can find models, validators and the full specification of the API (it can be accessed via [localhost:3000/swagger.json](http://localhost:3000/swagger.json)).

For convenience, this also includes a UI that represents the specs of the API and also includes an in-browser client: [localhost:3000/swaggerui/](http://localhost:3000/swaggerui/**

## Other considerations

*Vendoring*

This project is using go modules with vendoring. You can clone it anywhere in your filesystem.

*Storage layer*

The interface [`PaymentsStore`](./internal/store/store.go) allows to implement different storage layers for payments: e.g. a database store.

For simplicity, this uses an [in-memory store](./internal/store/volatile.go**. This is convenient because it is easy to work with, simple to use in tests and fast.

*Validations*

The service does some model-level validations. More strict validations could be defined in the swagger specs, and they would be used automatically.
