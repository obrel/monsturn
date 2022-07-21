package pdb

import (
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/monsturn/config"
)

var (
	db *sqlx.DB
)

func Init(cfg config.Db) error {
	d, err := url.Parse(cfg.Dsn)
	if err != nil {
		log.For("pdb", "init").Error(err)
		return err
	}

	d.User = url.UserPassword(d.User.Username(), "-FILTERED-")
	log.For("pdb", "init").Infof("Opening database on %s", d.String())

	db, err = sqlx.Open("postgres", cfg.Dsn)
	if err != nil {
		log.For("pdb", "init").Error(err)
		return err
	}

	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxLifetime))

	return nil
}

func tx(f func(*sqlx.Tx) error) error {
	t := db.MustBegin()

	err := f(t)
	if err != nil {
		log.For("pdb", "tx").Error(err)
		t.Rollback()
		return err
	}

	t.Commit()
	return nil
}
