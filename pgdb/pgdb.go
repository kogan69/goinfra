package pgdb

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/kogan69/goinfra/utils"
	"github.com/pkg/errors"
)

type PgDb struct {
	pool *pgxpool.Pool
}

type TxPg struct {
	tx pgx.Tx
	context.Context
}

func NewPgDbWithLog(dbUrl, logLevel string) (*PgDb, error) {
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	pgLogger := NewPgLogger(logLevel)
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgLogger,
		LogLevel: 0,
	}

	config.AfterConnect = TypeRegister

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &PgDb{pool: conn}, nil
}

func NewPgLogger(logLevel string) *PgLogger {
	pgLevel, err := tracelog.LogLevelFromString(logLevel)
	if err != nil {
		slog.Log(context.Background(),
			slog.LevelError,
			fmt.Sprintf("NewPgDbWithLog: failed to parse the logLevel %s with error: %s. defaulting to %s",
				logLevel,
				err.Error(),
				slog.LevelError.String()))
	}

	var level slog.Level
	switch pgLevel {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		level = slog.LevelDebug
	case tracelog.LogLevelInfo:
		level = slog.LevelInfo
	case tracelog.LogLevelWarn:
		level = slog.LevelWarn
	case tracelog.LogLevelError:
		level = slog.LevelError
	default:
		level = slog.LevelError
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	return &PgLogger{
		logger: logger,
		level:  level,
	}
}

type PgLogger struct {
	logger *slog.Logger
	level  slog.Level
}

func (p PgLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	if level == tracelog.LogLevelNone {
		return
	}
	attrs := make([]slog.Attr, 0)
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}
	p.logger.LogAttrs(ctx, p.level, msg, attrs...)
}

func (p *PgDb) Begin(ctx context.Context) (*TxPg, error) {

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		err = errors.Wrap(err, utils.FunctionName())
		return nil, err
	}
	txPg := &TxPg{
		tx:      tx,
		Context: ctx,
	}
	return txPg, err
}

func (p *PgDb) Query(ctx context.Context, sql string, params ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, params...)
}

func (p *PgDb) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *PgDb) ScanRow(dest any, rows pgx.Rows) error {
	return pgxscan.ScanRow(dest, rows)
}
func (p *PgDb) ScanOne(dst interface{}, rows pgx.Rows) error {
	return pgxscan.ScanOne(dst, rows)
}
func (p *PgDb) ScanAll(dst interface{}, rows pgx.Rows) error {
	return pgxscan.ScanAll(dst, rows)
}
func (p *PgDb) CloseRows(rows pgx.Rows) {
	if rows != nil {
		rows.Close()
	}
}

func (t *TxPg) Commit() error {
	return t.tx.Commit(t.Context)
}

func (t *TxPg) Rollback() error {
	return t.tx.Commit(t.Context)
}

func (t *TxPg) Tx() pgx.Tx {
	return t.tx
}
