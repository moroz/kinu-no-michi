-- +goose Up
-- +goose StatementBegin
create table orders (
  id uuid primary key,
  email_encrypted bytea not null,
  inserted_at timestamp(0) not null default (now() at time zone 'utc'),
  updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders;
-- +goose StatementEnd
