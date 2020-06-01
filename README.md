# retail-sample

An attempt to design a stock management system with domain driven principles in mind.
This was achieved by writing both the front and back-end code in such a way that their 'business' related logic is separated from the delivery and persistence code.
This allows to easily swap frameworks, libraries and databases without interfering with the business rules.

## Quick start

Given that you have `yarn` and `go` installed, just run:

```sh
make clean build run/mem
```

To run the app in docker compose with a postgres instance:

```sh
make clean build run/docker
```

To run the acceptance tests against a running instance - enter the `test/acceptance` folder then run:

```sh
yarn test
```

## Summary

The app has a Go backend exposing a REST interface and is using postgres for persistence.

The frontend is a SPA (single page application) written as a single static HTML page using Bootstrap.

Both front and back-end are written in such a way that their 'business' related code is separated from the delivery code:

- for the back-end - all the implementation details are placed in the `cmd` folder, while all the business rules are in the `internal` directory.

- for the front-end - the `web/src/app` folder contains all the framework free code - the domain objects, rest clients and page behaviour (error messages, controls animation) while `web/src/jquery` or `web/src/plan` (WIP) contains the means by which this logic in bound to the web page components.

## Development

### back-end

> make run

Linter is going to be ran as part of building the runnable docker image:

> make build/docker

### front-end

enter `web/` folder

> make watch

## Business domain explained

### inventory

- a list of unique item names

### stock

- a list of items that are present in the inventory AND have been provisioned

### recipe

- a list of recipes, each recipe has a unique name and a list of ingredients with their quantities

### order

- a list of orders that have been successfully placed.
- if the order does not have enough stock, it will be rejected
- if there is enough stock, the corresponding values will be substracted from the stock

> to be continued...

## TODO

### tech

- go api tests
- errors in jsonapi format
- document packages with godoc
- pagination
- linter yml
- pprof
- auth and roles

### business

- view provision log
- cancel a provision
- cancel an order
- ...
