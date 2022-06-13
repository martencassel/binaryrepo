package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/martencassel/binaryrepo"
	"github.com/pkg/errors"
)

type Database struct {
	*pgxpool.Pool
	name string
	cfg *binaryrepo.Config
}

func (t *txnSession) doQuery(q sq.SelectBuilder) (pgx.Rows, error) {
	sqlStr, args, err := q.ToSql()
	log.Println(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("query to sql: %v", err)
	}
	return t.Query(t.ctx, sqlStr, args...)
}


func createPool(ctx context.Context, cfg binaryrepo.Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	pgxCfg.MaxConns = cfg.MaxOpenConnections
	pgxCfg.MaxConns = cfg.MaxIdleConnections
	pgxCfg.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetime)
	pool, err := pgxpool.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func createDBPool(ctx context.Context, cfg binaryrepo.Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	pgxCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Printf("Registering types: repotype\n")
		repoType, err := pgxtype.LoadDataType(ctx, conn, conn.ConnInfo(), "repotype")
		log.Println(repoType, err)
		if err != nil {
			return err
		}
		pkgType, err := pgxtype.LoadDataType(ctx, conn, conn.ConnInfo(), "pkgtype")
		log.Println(pkgType, err)
		if err != nil {
			return err
		}
		conn.ConnInfo().RegisterDataType(repoType)
		conn.ConnInfo().RegisterDataType(pkgType)
		return nil
	}
	pgxCfg.MaxConns = cfg.MaxOpenConnections
	pgxCfg.MaxConns = cfg.MaxIdleConnections
	pgxCfg.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetime)
	pool, err := pgxpool.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func migrateDB(ctx context.Context, pool *pgxpool.Pool) error {
	pgxPoolConn, err := pool.Acquire(ctx)
	if err != nil {
			return errors.Wrap(err, "failed to acquire connection")
	}
	defer pgxPoolConn.Release()
	 var pgxConn *pgx.Conn = pgxPoolConn.Conn()
	migrator, err := migrate.NewMigrator(ctx, pgxConn, MigrationTable)
	if err != nil {
			return errors.Wrap(err, "failed to create migrator")
	}
	for _, m := range schemaMigrations {
			migrator.AppendMigration(m.Name, m.UpSQL, m.DownSQL)
	}
	if err := migrator.Migrate(ctx); err != nil {
			return errors.Wrap(err, "failed to migrate")
	}
	return nil
}


func dbCreateStatement(dbName string) string {
	return fmt.Sprintf("CREATE DATABASE %s", dbName)
}

func dbExistsQuery(dbName string) string {
	return fmt.Sprintf("SELECT true FROM pg_database WHERE datname='%s'", dbName)
}

// Open connects to the database and creates the database if its missing
func Open(ctx context.Context, cfg binaryrepo.Config) (*Database, error) {
	db := Database {
		name: cfg.DBName,
		cfg: &cfg,
	}
	if cfg.CreateDB && cfg.DBName != "postgres" {
		log.Printf("Creating database: %s", cfg.DBName)
		tmpCfg := cfg
		tmpCfg.DBName = "postgres"
		tmpPg, err := createPool(ctx, tmpCfg)
		log.Println("After createPool")
		defer tmpPg.Close()
		log.Println(tmpPg)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var exists bool
		if err := tmpPg.QueryRow(ctx, dbExistsQuery(cfg.DBName)).Scan(&exists); err != nil {
			log.Println(err)
			if  err != pgx.ErrNoRows {
				log.Fatalf("Failed to check if database exists: %v", err)
				return nil, err
			}
			log.Printf("Creating databases with %s", dbCreateStatement(cfg.DBName))
			log.Printf("%s", dbCreateStatement(cfg.DBName))
			_, err := tmpPg.Exec(ctx, dbCreateStatement(cfg.DBName))
			if err != nil {
				log.Fatalf("Failed to create database: %v", err)
				return nil, errors.Wrapf(err, "failed to create database: %s", cfg.DBName)
			}
		}
	}
	pool, err := createPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	cfg = *db.cfg
	err = migrateDB(ctx, pool)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, errors.Wrapf(err, "failed to migrate database: %s", cfg.DBName)
	}

	dbPool, err := createDBPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	db.Pool = dbPool
	return &db, nil
}

func InsertData(pool *pgxpool.Pool, name string, email string) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, "password")
	if err != nil {
		log.Fatalf("Failed to insert row: %v", err)
	}
}

func QueryData(pool *pgxpool.Pool) {
	ctx := context.Background()
	rows, err := pool.Query(ctx, "SELECT email FROM users")
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("Email: %s\n", email)
	}
	var email string
	pool.QueryRow(ctx, "SELECT email FROM users WHERE id = $1", 6).Scan(&email)
	fmt.Printf("Email for id = 1 is Email: %s\n", email)
}



func build() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
