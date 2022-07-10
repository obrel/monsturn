package mdb

import (
	"database/sql"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/monsturn/config"
)

var (
	db *sql.DB
)

func Init(cfg config.Db) error {
	d, err := url.Parse(cfg.Dsn)
	if err != nil {
		log.For("mdb", "init").Error(err)
		return err
	}

	d.User = url.UserPassword(d.User.Username(), "-FILTERED-")
	log.For("mdb", "init").Infof("Opening database on %s", d.String())

	db, err = sql.Open("mysql", cfg.Dsn)
	if err != nil {
		log.For("mdb", "init").Error(err)
		return err
	}

	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxLifetime))

	return nil
}
