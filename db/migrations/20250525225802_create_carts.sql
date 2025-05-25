-- +goose Up
-- +goose StatementBegin
create table carts (
  id uuid primary key,
  inserted_at timestamp(0) default (now() at time zone 'utc'),
  updated_at timestamp(0) default (now() at time zone 'utc')
);

create table cart_items (
  id uuid primary key,
  cart_id uuid not null references carts (id) on delete cascade,
  product_id uuid not null references products (id) on delete cascade,
  quantity decimal not null default 1,
  inserted_at timestamp(0) default (now() at time zone 'utc'),
  updated_at timestamp(0) default (now() at time zone 'utc'),
  unique (cart_id, product_id),
  check (quantity > 0)
);

create index on cart_items (cart_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cart_items;
drop table carts;
-- +goose StatementEnd
