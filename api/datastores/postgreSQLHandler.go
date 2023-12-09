package datastores

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rigved-desai/paryatan-backend/api/interfaces"
)

type PostgreSQLHanlder struct {
	ConnPool *pgxpool.Pool
}

func (db *PostgreSQLHanlder) Query(query string, args ...interface{}) (interfaces.DBRow, error) {
	rows, err := db.ConnPool.Query(context.Background(), query, args...)
	if err != nil {
		return &PGXRow{}, err
	}
	return &PGXRow{rows: rows}, nil
}

func (db *PostgreSQLHanlder) Execute(query string, args ...interface{}) (pgx.Row, error) {
	row := db.ConnPool.QueryRow(context.Background(), query, args...)
	return row, nil
}

type PGXRow struct {
	rows pgx.Rows
}

func (r *PGXRow) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *PGXRow) Next() bool {
	return r.rows.Next()
}