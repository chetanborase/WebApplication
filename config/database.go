package config

import (
	"WebApplication/env"
	"WebApplication/internal/logger"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	db   *sqlx.DB
	ones = &sync.Once{}
)

type dbConfig struct {
	driverName string
	user       string
	pass       string
	database   string
	address    string
}

func newConf() dbConfig {
	return dbConfig{
		driverName: env.Get(env.DbDriveName, "mysql"),
		user:       env.Get(env.DbUserName, ""),
		pass:       env.Get(env.DbPassword, ""),
		database:   env.Get(env.DbSchema, ""),
		address:    env.Get(env.DbServerAddress, ":8080"),
	}
}
func GetDB() *sqlx.DB {
	ones.Do(
		func() {
			env.LoadProfile()
			var err error
			conf := newConf()
			db, err = sqlx.Open(conf.driverName, conf.getSourceString())
			if err != nil {
				log.Fatal(err)
			}

			m, err := migrate.New("file://sql",
				fmt.Sprintf("%s://%s", conf.driverName, conf.getSourceString()))
			if err != nil {
				logger.Error(err)
			}
			if err = m.Up(); err != nil {
				if !strings.Contains(err.Error(), "no change") {
					logger.Error(err)
				}
			}
			db.SetConnMaxIdleTime(10 * time.Second)
		})
	return db
}
func (conf dbConfig) getSourceString() string {
	flag := "?parseTime=true"
	return fmt.Sprintf("%s:%s@tcp(%s)/%s%s",
		conf.user, conf.pass, conf.address, conf.database, flag)
}
