package binaryrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)


type Config struct {
	// Database host, required
	Host string
	// Database port, required
	Port string
	// Database name, required
	DBName string
	// Database user, required
	User string
	// Database password, required
	Password string
	// Limit on open connections to the database, required
	MaxOpenConnections int32
	// Limit on the idle connections to the database, required
	MaxIdleConnections int32
	// Connection max life time, optional
	ConnMaxLifetime int32
	// Auto create database if not exists, optional
	CreateDB bool
}


const (
	ConnTimeoutSec = 10
)

func (c *Config) ConnectionString() string {
	return strings.Join([]string{
		"host=" + c.Host,
		"port=" + c.Port,
		"dbname=" + c.DBName,
		"user=" + c.User,
		"password=" + c.Password,
		fmt.Sprintf("connect_timeout=%d", ConnTimeoutSec),
	}, " ")
}



type Database interface {
	ConnectionString() string;
	Open(ctx context.Context, cfg Config) (*Database, error)
	InsertData(pool *pgxpool.Pool, name string, email string)
	QueryData(pool *pgxpool.Pool)

}