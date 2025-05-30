// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: products.sql

package queries

import (
	"context"
)

const listProducts = `-- name: ListProducts :many
select id, title, slug, base_price_eur, description, image_url, inserted_at, updated_at from products order by id desc
`

func (q *Queries) ListProducts(ctx context.Context) ([]*Product, error) {
	rows, err := q.db.Query(ctx, listProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Slug,
			&i.BasePriceEur,
			&i.Description,
			&i.ImageUrl,
			&i.InsertedAt,
			&i.UpdatedAt,
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
