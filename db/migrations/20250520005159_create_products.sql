-- +goose Up
-- +goose StatementBegin
create table products (
  id uuid primary key,
  title text not null,
  slug text not null unique,
  base_price_eur decimal not null default 0,
  description text,
  inserted_at timestamp(0) not null default (now() at time zone 'utc'),
  updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table products;
-- +goose StatementEnd
