package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sudhir512kj/ecommerce_backend/config"
)

type postgresDatabase struct {
	Db *sql.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic("failed to connect database")
		}

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *sql.DB {
	return p.Db
}
