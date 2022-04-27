package binaryrepodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

const (
	ConnTimeoutSec = 10
	MigrationTable = "binaryrepo_schema_migrations"
)

// Config expected to be used.
type Config struct {
	// Database host, required
	Host string
	// Database port, required
	Port int
	// Database name, required
	DBName string
	// Database user, required
	User string
	// Database password, required
	Password string
	// Limit the open connections in the pool, required
	MaxOpenConnections int32
	// Limit on the idle connections in the pool, required
	MaxIdleConnections int32
	// Limit on the lifetime of each connection in the pool, required
	ConnMaxLifetime time.Duration
	// Automatically create database if it doesn't exist
	CreateDB bool
}

// Database provides postgres connection pool
type Database struct {
	pool *pgxpool.Pool
}

// ConnString build a connection string from Config
func (c *Config) ConnString() string {
	return strings.Join(
		append([]string{
			fmt.Sprintf("host=%s", c.Host),
			fmt.Sprintf("port=%d", c.Port),
			fmt.Sprintf("user=%s", c.User),
			fmt.Sprintf("password=%s", c.Password),
			fmt.Sprintf("dbname=%s", c.DBName),
			fmt.Sprintf("timezone=%s", "utc"),
			fmt.Sprintf("connect_timeout=%d", ConnTimeoutSec)}), " ")
}

// Open connects to the database and creates the database if it's missing
func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	if cfg.CreateDB && cfg.DBName != "postgres" {
		log.Println("Creating database...")

		tempCfg := cfg
		tempCfg.DBName = "postgres"
		tempPgPool, err := createPool(ctx, tempCfg)
		if err != nil {
			return nil, err
		}
		sqlExistsDbQuery := fmt.Sprintf("SELECT true FROM pg_database WHERE datname='%s'", cfg.DBName)
		log.Println(sqlExistsDbQuery)
		var databaseExists bool
		if err := tempPgPool.QueryRow(ctx, sqlExistsDbQuery).Scan(&databaseExists); err != nil {
			if err != pgx.ErrNoRows {
				return nil, err
			}
			sqlCreateDatabaseStmt := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
			_, err := tempPgPool.Exec(ctx, sqlCreateDatabaseStmt)
			if err != nil {
				return nil, fmt.Errorf("failed to create postgres database: %s", cfg.DBName)
			}
		}
	}
	pool, err := createPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(ctx, pool); err != nil {
		return nil, err
	}

	return pool, nil
}

// OpenOrFail connects to the database and creates the database if it's missing,
// it throws an fatal error on failure
func OpenOrFail(ctx context.Context, cfg Config) (*pgxpool.Pool) {
	pg, err := Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	return pg
}

// createPool creates a postgres connection pool for a given configuration
func createPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
		// Create pool
		pgxCfg, err := pgxpool.ParseConfig(cfg.ConnString())
		if err != nil {
			return nil, err
		}
		pgxCfg.MaxConns = cfg.MaxOpenConnections
		pgxCfg.MinConns = cfg.MaxIdleConnections
		pgxCfg.MaxConnLifetime = cfg.ConnMaxLifetime

		pool, err := pgxpool.ConnectConfig(ctx, pgxCfg)
		if err != nil {
			return nil, err
		}
		return pool, nil
}

func migrateDB(ctx context.Context, pool *pgxpool.Pool) error {
	pgxconn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer pgxconn.Release()
	migrator, err := newMigrator(ctx, pgxconn.Conn())
	if err != nil {
		return err
	}
	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to migrate database: %s", err)
	}
	return nil
}

func newMigrator(ctx context.Context, conn *pgx.Conn) (*migrate.Migrator, error) {
	migrator, err := migrate.NewMigrator(ctx, conn, MigrationTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %s", err)
	}
	for _, m := range schemaMigrations {
		migrator.AppendMigration(m.Name, m.UpSQL, m.DownSQL)
	}
	return migrator, nil
}

/*	ctx := context.Background()
	cfg := Config{
		Host: "localhost",
		Port: 5432,
		DBName: "binaryrepo",
		User: "postgres",
		Password: "postgres",
		MaxOpenConnections: 10,
		MaxIdleConnections: 10,
		ConnMaxLifetime: time.Hour,
		CreateDB: true,
	}
	pool, err  := Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println(pool)
*/

