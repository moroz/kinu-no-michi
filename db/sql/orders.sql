-- name: InsertOrder :one
insert into orders (id, email_encrypted, grand_total_eur, grand_total_btc, exchange_rate) values ($1, $2, $3, $4, $5) returning *;
