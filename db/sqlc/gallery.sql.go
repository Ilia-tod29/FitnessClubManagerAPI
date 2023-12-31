// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: gallery.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createGalleryItem = `-- name: CreateGalleryItem :one
INSERT INTO gallery (
    image
) VALUES (
    $1
) RETURNING id, image, created_at
`

func (q *Queries) CreateGalleryItem(ctx context.Context, image pgtype.Text) (Gallery, error) {
	row := q.db.QueryRow(ctx, createGalleryItem, image)
	var i Gallery
	err := row.Scan(&i.ID, &i.Image, &i.CreatedAt)
	return i, err
}

const deleteGalleryItem = `-- name: DeleteGalleryItem :one
DELETE FROM gallery
WHERE id = $1
RETURNING id, image, created_at
`

func (q *Queries) DeleteGalleryItem(ctx context.Context, id int64) (Gallery, error) {
	row := q.db.QueryRow(ctx, deleteGalleryItem, id)
	var i Gallery
	err := row.Scan(&i.ID, &i.Image, &i.CreatedAt)
	return i, err
}

const getGalleryItem = `-- name: GetGalleryItem :one
SELECT id, image, created_at FROM gallery
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGalleryItem(ctx context.Context, id int64) (Gallery, error) {
	row := q.db.QueryRow(ctx, getGalleryItem, id)
	var i Gallery
	err := row.Scan(&i.ID, &i.Image, &i.CreatedAt)
	return i, err
}

const listAllGalleryItems = `-- name: ListAllGalleryItems :many
SELECT id, image, created_at FROM gallery
ORDER BY id
`

func (q *Queries) ListAllGalleryItems(ctx context.Context) ([]Gallery, error) {
	rows, err := q.db.Query(ctx, listAllGalleryItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Gallery{}
	for rows.Next() {
		var i Gallery
		if err := rows.Scan(&i.ID, &i.Image, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGalleryItems = `-- name: ListGalleryItems :many
SELECT id, image, created_at FROM gallery
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListGalleryItemsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListGalleryItems(ctx context.Context, arg ListGalleryItemsParams) ([]Gallery, error) {
	rows, err := q.db.Query(ctx, listGalleryItems, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Gallery{}
	for rows.Next() {
		var i Gallery
		if err := rows.Scan(&i.ID, &i.Image, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
