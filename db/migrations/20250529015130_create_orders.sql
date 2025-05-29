-- +goose Up
-- +goose StatementBegin
create table orders (
  id uuid primary key,
  email_encrypted bytea not null,
  grand_total_eur decimal not null,
  grand_total_btc decimal not null,
  exchange_rate decimal not null,
  inserted_at timestamp(0) not null default (now() at time zone 'utc'),
  updated_at timestamp(0) not null default (now() at time zone 'utc')
);

create table order_line_items (
  id uuid primary key,
  order_id uuid not null references orders (id) on delete cascade,
  product_id uuid references products (id) on delete set null,
  quantity decimal not null,
  unit_price_eur decimal not null,
  subtotal_btc decimal not null,
  product_title text not null,
  inserted_at timestamp(0) not null default (now() at time zone 'utc'),
  updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_line_items;
drop table orders;
-- +goose StatementEnd
