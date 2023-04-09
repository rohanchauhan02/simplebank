// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: entry.sql

package db

import (
	"context"
)

const createEntries = `-- name: CreateEntries :one
INSERT INTO entries (
    account_id,
    account
) VALUES (
  $1, $2
) RETURNING id, account_id, account, created_at
`

type CreateEntriesParams struct {
	AccountID int64 `json:"accountID"`
	Account   int64 `json:"account"`
}

func (q *Queries) CreateEntries(ctx context.Context, arg CreateEntriesParams) (Entry, error) {
	row := q.queryRow(ctx, q.createEntriesStmt, createEntries, arg.AccountID, arg.Account)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Account,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entries
where id = $1
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteEntryStmt, deleteEntry, id)
	return err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, account, created_at FROM entries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.queryRow(ctx, q.getEntryStmt, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Account,
		&i.CreatedAt,
	)
	return i, err
}

const listEntry = `-- name: ListEntry :many
SELECT id, account_id, account, created_at FROM entries
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListEntryParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEntry(ctx context.Context, arg ListEntryParams) ([]Entry, error) {
	rows, err := q.query(ctx, q.listEntryStmt, listEntry, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Account,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEntry = `-- name: UpdateEntry :exec
UPDATE entries
SET account_id = $1
where id = $2
`

type UpdateEntryParams struct {
	AccountID int64 `json:"accountID"`
	ID        int64 `json:"id"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) error {
	_, err := q.exec(ctx, q.updateEntryStmt, updateEntry, arg.AccountID, arg.ID)
	return err
}
