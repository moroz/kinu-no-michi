-- name: InsertOrder :one
insert into orders (id, email_encrypted) values ($1, $2) returning *;
