package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgre struct {
	DB *pgxpool.Pool
}

func NewDB() (*Postgre, error) {
	dbPool, err := pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		return nil, err
	}
	postgre := &Postgre{
		DB: dbPool,
	}
	return postgre, nil
}
