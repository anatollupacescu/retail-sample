create table inventory (
  id serial primary key,
  name varchar(32) not null
);

create table stock (
  id serial primary key,
  itemID int not null,
  qty int not null
);