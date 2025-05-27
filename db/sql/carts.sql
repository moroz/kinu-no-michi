-- name: InsertCart :exec
insert into carts (id) values ($1);

-- name: InsertCartItem :one
insert into cart_items as ci (id, cart_id, product_id, quantity)
values ($1, $2, $3, $4)
on conflict (cart_id, product_id) do update
set quantity = ci.quantity + excluded.quantity,
updated_at = now() at time zone 'utc'
returning *;

-- name: GetCartByID :one
select c.id, count(ci.id) item_count, sum(ci.quantity * p.base_price_eur)::decimal grand_total
from carts c
join cart_items ci on c.id = ci.cart_id
join products p on p.id = ci.product_id
where c.id = $1 group by 1 limit 1;
