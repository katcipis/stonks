# User Manager

This is a service responsible for managing users, providing a way
to create, list and delete users through a simple to use API.

# API

The reference documentation for the API can be found [here](docs/api.md).

# Development

## Dependencies

To run the automation provided here you need to have installed:

* [Make](https://www.gnu.org/software/make/)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

If you fancy building things (or running tests) with no extra layers
you can install the latest [Go](https://golang.org/doc/install). It can
be useful to quickly run isolated tests like this:

```sh
go test -race ./...
```

But for running the integration tests it gets challenging since the
project depends on two different databases, [Redis](https://redis.io/)
and [Postgres](https://www.postgresql.org/).

It is recommended to just use the provided automation through make,
it will help you achieve consistent results.


## Running tests

To run isolated tests just run:

```
make test
```

And for integration tests:

```
make test-integration
```

To check locally the coverage from all the tests run:

```
make coverage
```

And it should open the coverage analysis in your browser.


## Linting

To lint code, just run:

```
make lint
```

## Releasing

To create an image ready for production deployment just run:

```
make image
```

And one will be created tagged with the git short revision of the
code used to build it, you can also specify an explicit version
if you want:

```
make image version=1.12.0
```
