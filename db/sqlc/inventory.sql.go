// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: inventory.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createInventoryItem = `-- name: CreateInventoryItem :one
INSERT INTO inventory (
    name,
    image
) VALUES (
    $1, $2
) RETURNING id, name, image, created_at
`

type CreateInventoryItemParams struct {
	Name  string      `json:"name"`
	Image pgtype.Text `json:"image"`
}

func (q *Queries) CreateInventoryItem(ctx context.Context, arg CreateInventoryItemParams) (Inventory, error) {
	row := q.db.QueryRow(ctx, createInventoryItem, arg.Name, arg.Image)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const deleteInventoryItem = `-- name: DeleteInventoryItem :one
DELETE FROM inventory
WHERE id = $1
RETURNING id, name, image, created_at
`

func (q *Queries) DeleteInventoryItem(ctx context.Context, id int64) (Inventory, error) {
	row := q.db.QueryRow(ctx, deleteInventoryItem, id)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const getInventoryItem = `-- name: GetInventoryItem :one
SELECT id, name, image, created_at FROM inventory
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetInventoryItem(ctx context.Context, id int64) (Inventory, error) {
	row := q.db.QueryRow(ctx, getInventoryItem, id)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const listAllInventoryItems = `-- name: ListAllInventoryItems :many
SELECT id, name, image, created_at FROM inventory
ORDER BY id
`

func (q *Queries) ListAllInventoryItems(ctx context.Context) ([]Inventory, error) {
	rows, err := q.db.Query(ctx, listAllInventoryItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Inventory{}
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listInventoryItems = `-- name: ListInventoryItems :many
SELECT id, name, image, created_at FROM inventory
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListInventoryItemsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListInventoryItems(ctx context.Context, arg ListInventoryItemsParams) ([]Inventory, error) {
	rows, err := q.db.Query(ctx, listInventoryItems, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Inventory{}
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateInventoryItem = `-- name: UpdateInventoryItem :one
UPDATE inventory
set name = $2,
    image = $3
WHERE id = $1
RETURNING id, name, image, created_at
`

type UpdateInventoryItemParams struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Image pgtype.Text `json:"image"`
}

func (q *Queries) UpdateInventoryItem(ctx context.Context, arg UpdateInventoryItemParams) (Inventory, error) {
	row := q.db.QueryRow(ctx, updateInventoryItem, arg.ID, arg.Name, arg.Image)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const updateInventoryItemImage = `-- name: UpdateInventoryItemImage :one
UPDATE inventory
set image = $2
WHERE id = $1
RETURNING id, name, image, created_at
`

type UpdateInventoryItemImageParams struct {
	ID    int64       `json:"id"`
	Image pgtype.Text `json:"image"`
}

func (q *Queries) UpdateInventoryItemImage(ctx context.Context, arg UpdateInventoryItemImageParams) (Inventory, error) {
	row := q.db.QueryRow(ctx, updateInventoryItemImage, arg.ID, arg.Image)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const updateInventoryItemName = `-- name: UpdateInventoryItemName :one
UPDATE inventory
set name = $2
WHERE id = $1
RETURNING id, name, image, created_at
`

type UpdateInventoryItemNameParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateInventoryItemName(ctx context.Context, arg UpdateInventoryItemNameParams) (Inventory, error) {
	row := q.db.QueryRow(ctx, updateInventoryItemName, arg.ID, arg.Name)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}
