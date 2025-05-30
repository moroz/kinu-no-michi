// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: carts.sql

package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const deleteCart = `-- name: DeleteCart :exec
delete from carts where id = $1
`

func (q *Queries) DeleteCart(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCart, id)
	return err
}

const getCartByID = `-- name: GetCartByID :one
select c.id, count(ci.id) item_count, sum(ci.quantity * p.base_price_eur)::decimal grand_total
from carts c
join cart_items ci on c.id = ci.cart_id
join products p on p.id = ci.product_id
where c.id = $1 group by 1 limit 1
`

type GetCartByIDRow struct {
	ID         uuid.UUID
	ItemCount  int64
	GrandTotal decimal.Decimal
}

func (q *Queries) GetCartByID(ctx context.Context, id uuid.UUID) (*GetCartByIDRow, error) {
	row := q.db.QueryRow(ctx, getCartByID, id)
	var i GetCartByIDRow
	err := row.Scan(&i.ID, &i.ItemCount, &i.GrandTotal)
	return &i, err
}

const getCartItemsByCartID = `-- name: GetCartItemsByCartID :many
select ci.id, ci.product_id, ci.quantity, p.base_price_eur, p.title
from cart_items ci
join products p on ci.product_id = p.id
where ci.cart_id = $1
`

type GetCartItemsByCartIDRow struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	Quantity     decimal.Decimal
	BasePriceEur decimal.Decimal
	Title        string
}

func (q *Queries) GetCartItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]*GetCartItemsByCartIDRow, error) {
	rows, err := q.db.Query(ctx, getCartItemsByCartID, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetCartItemsByCartIDRow
	for rows.Next() {
		var i GetCartItemsByCartIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Quantity,
			&i.BasePriceEur,
			&i.Title,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCart = `-- name: InsertCart :exec
insert into carts (id) values ($1)
`

func (q *Queries) InsertCart(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, insertCart, id)
	return err
}

const insertCartItem = `-- name: InsertCartItem :one
insert into cart_items as ci (id, cart_id, product_id, quantity)
values ($1, $2, $3, $4)
on conflict (cart_id, product_id) do update
set quantity = ci.quantity + excluded.quantity,
updated_at = now() at time zone 'utc'
returning id, cart_id, product_id, quantity, inserted_at, updated_at
`

type InsertCartItemParams struct {
	ID        uuid.UUID
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  decimal.Decimal
}

func (q *Queries) InsertCartItem(ctx context.Context, arg InsertCartItemParams) (*CartItem, error) {
	row := q.db.QueryRow(ctx, insertCartItem,
		arg.ID,
		arg.CartID,
		arg.ProductID,
		arg.Quantity,
	)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.ProductID,
		&i.Quantity,
		&i.InsertedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
