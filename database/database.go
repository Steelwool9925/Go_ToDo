package database

import (
	cfg "Go_Test/config"
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module exports the DB connection provider for FX.
var Module = fx.Options(
	fx.Provide(NewDBConnection),
)

type DBConnectionParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *cfg.Config
	Logger    *zap.Logger
}

// NewDBConnection creates and manages a database connection lifecycle.
func NewDBConnection(p DBConnectionParams) (*sql.DB, error) {
	p.Logger.Info("Attempting to connect to database", zap.String("db_host", p.Config.DBHost), zap.String("db_name", p.Config.DBName))

	db, err := sql.Open("mysql", p.Config.DBDSN)
	if err != nil {
		p.Logger.Error("Failed to open database connection", zap.Error(err))
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("Pinging database...")
			if err := db.PingContext(ctx); err != nil {
				p.Logger.Error("Failed to ping database on start", zap.Error(err))
				db.Close()
				return err
			}
			p.Logger.Info("Database connection established and pinged successfully.")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Closing database connection")
			if err := db.Close(); err != nil {
				p.Logger.Error("Failed to close database connection", zap.Error(err))
				return err
			}
			p.Logger.Info("Database connection closed.")
			return nil
		},
	})

	return db, nil
}
