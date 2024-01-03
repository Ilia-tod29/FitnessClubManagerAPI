// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: subscription.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO subscriptions (
    user_id,
    start_date,
    end_date
) VALUES (
    $1, $2, $3
) RETURNING id, user_id, start_date, end_date, created_at
`

type CreateSubscriptionParams struct {
	UserID    int64       `json:"user_id"`
	StartDate pgtype.Date `json:"start_date"`
	EndDate   pgtype.Date `json:"end_date"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRow(ctx, createSubscription, arg.UserID, arg.StartDate, arg.EndDate)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
	)
	return i, err
}

const getSubscription = `-- name: GetSubscription :one
SELECT id, user_id, start_date, end_date, created_at FROM subscriptions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSubscription(ctx context.Context, id int64) (Subscription, error) {
	row := q.db.QueryRow(ctx, getSubscription, id)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
	)
	return i, err
}

const listAllSubscriptions = `-- name: ListAllSubscriptions :many
SELECT id, user_id, start_date, end_date, created_at FROM subscriptions
ORDER BY id
`

func (q *Queries) ListAllSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := q.db.Query(ctx, listAllSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscription{}
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StartDate,
			&i.EndDate,
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

const listSubscriptions = `-- name: ListSubscriptions :many
SELECT id, user_id, start_date, end_date, created_at FROM subscriptions
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListSubscriptionsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListSubscriptions(ctx context.Context, arg ListSubscriptionsParams) ([]Subscription, error) {
	rows, err := q.db.Query(ctx, listSubscriptions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscription{}
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StartDate,
			&i.EndDate,
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
