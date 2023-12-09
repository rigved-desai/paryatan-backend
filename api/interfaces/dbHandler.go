package interfaces

import "github.com/jackc/pgx/v4"

// will be implemented by postgreSQLHandler in datastores pkg
type DBHandler interface {
    Query(query string, args ...interface{}) (DBRow, error)
    Execute(query string, args ...interface{}) (pgx.Row, error)
}

type DBRow interface {
    Scan(dest ...interface{}) error
    Next() bool
}