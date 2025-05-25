-- name: InsertCart :exec
insert into carts (id) values ($1);

-- name: InsertCartItem :one
insert into cart_items as ci (id, cart_id, product_id, quantity)
values ($1, $2, $3, $4)
on conflict (cart_id, product_id) do update
set quantity = ci.quantity + excluded.quantity,
updated_at = now() at time zone 'utc'
returning *;
