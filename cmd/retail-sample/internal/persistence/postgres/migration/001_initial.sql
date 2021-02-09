CREATE TABLE inventory (
  id bigserial PRIMARY KEY,
  name varchar(36) NOT NULL,
  enabled bool NOT NULL DEFAULT TRUE
);

CREATE UNIQUE INDEX idx_inventory_name ON inventory(name);

CREATE TABLE stock (
  id bigserial PRIMARY KEY,
  inventoryid bigint REFERENCES inventory (id) UNIQUE,
  quantity int NOT NULL
);

CREATE TABLE recipe (
  id bigserial PRIMARY KEY,
  name varchar(36) NOT NULL UNIQUE,
  enabled bool NOT NULL DEFAULT TRUE
);

CREATE TABLE recipe_ingredient (
  id bigserial PRIMARY KEY,
  recipeid bigint REFERENCES recipe (id),
  inventoryid bigint REFERENCES inventory (id),
  quantity smallint NOT NULL,
  UNIQUE (recipeid, inventoryid)
);

CREATE TABLE outbound_order (
  id bigserial PRIMARY KEY,
  recipeid int REFERENCES recipe (id) NOT NULL,
  quantity int NOT NULL,
  orderdate timestamp DEFAULT now(),
  UNIQUE (recipeid, orderdate)
);

CREATE TABLE provisionlog (
  id bigserial PRIMARY KEY,
  inventoryid bigint REFERENCES inventory (id),
  quantity smallint NOT NULL,
  provisiondate timestamp DEFAULT now()
);

