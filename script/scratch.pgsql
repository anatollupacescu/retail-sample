-- SELECT
--   r.id,
--   r.name,
--   i.id,
--   ri.quantity
-- FROM
--   recipe_ingredient ri,
--   recipe r,
--   inventory i
-- WHERE
--   ri.recipeid = r.id
--   AND ri.inventoryid = i.id;

DROP TABLE outbound_order;

DROP TABLE recipe_ingredient;

DROP TABLE recipe;

DROP TABLE stock;

DROP TABLE provisionlog;

DROP TABLE inventory;

CREATE TABLE inventory (
  id bigserial PRIMARY KEY,
  name varchar(36) NOT NULL
);

CREATE TABLE stock (
  id bigserial PRIMARY KEY,
  inventoryid bigint REFERENCES inventory (id) UNIQUE,
  quantity int NOT NULL
);

CREATE TABLE recipe (
  id bigserial PRIMARY KEY,
  name varchar(36) NOT NULL UNIQUE
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

