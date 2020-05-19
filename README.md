# retail-sample

**TODO**

godoc
inmemory as flag

disabling a certain inventory item
view recipe items
cancel a provision

**Summary**

An attempt to design a stock management system with domain driven principles in mind.

It has a Go backend exposing a REST interface and using postgres for persistence.

The frontend is a SPA (single page application) with a single static HTML page written using Bootstrap.

Both front and back-end are written in such a way that their 'business' related code is separated from the delivery code:

- for the back-end - all the implementation details are placed in the `cmd` folder, while all the business rules are in the `internal` directory.

- for the front-end - the `web/src/app` folder contains all the framework free code - the domain objects, rest clients and page behaviour (error messages, controls animation) while `web/src/jquery` or `web/src/plan` (WIP) contains the means by which this logic in bound to the web page components.

**Usage**

##front-end

enter web/ folder

make gen

make build

##back-end

make build

docker-compose up --build

to start the backend

> make run

to start frontend

> cd web; ./gen_static.sh; yarn start

The code is located in the `web` directory

## inventory:

- a list of unique item names

## stock:

- a list of items that are present in the inventory AND have been provisioned

## recipe:

- a list of recipes, each recipe has a unique name and a list of ingredients with their quantities

## order:

- a list of orders that have been successfully placed.
- if the order does not have enough stock, it will be rejected
- if there is enough stock, the corresponding values will be substracted from the stock
