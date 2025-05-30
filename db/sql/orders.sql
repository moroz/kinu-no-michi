-- name: InsertOrder :one
insert into orders (id, email_encrypted, grand_total_eur, grand_total_btc, exchange_rate) values ($1, $2, $3, $4, $5) returning *;

-- name: InsertOrderLineItem :exec
insert into order_line_items (id, order_id, product_id, quantity, product_unit_price, product_title) values ($1, $2, $3, $4, $5, $6);
