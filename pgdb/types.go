package pgdb

import (
	"context"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
)

func TypeRegister(_ context.Context, conn *pgx.Conn) (err error) {
	pgxuuid.Register(conn.TypeMap())
	pgxdecimal.Register(conn.TypeMap())
	return
}
