package rdb

import (
	"net/url"
	"strconv"
	"time"

	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/monsturn/config"

	"github.com/gomodule/redigo/redis"
)

var (
	db *redis.Pool
)

func Init(cfg config.Redis) error {
	var password string
	var database string

	d, err := url.Parse(cfg.Dsn)
	if err != nil {
		log.For("rdb", "init").Error(err)
		return err
	}

	if d.User != nil {
		if pwd, ok := d.User.Password(); ok {
			password = pwd
		} else {
			password = d.User.Username()
		}
	}

	if d.Path != "" {
		db, err := strconv.ParseInt(d.Path, 10, 64)
		if err == nil {
			database = strconv.Itoa(int(db))
		}
	}

	d.User = url.UserPassword(d.User.Username(), "-FILTERED-")
	log.For("rdb", "init").Infof("Opening redis on %s", d.String())

	db = &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		IdleTimeout: time.Duration(cfg.MaxTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", d.Host)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if database != "" {
				if _, err := c.Do("SELECT", database); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}

	_, err = db.Get().Do("PING")
	if err != nil {
		return err
	}

	return nil
}

func GetPool() *redis.Pool {
	return db
}

func GetConn() redis.Conn {
	return db.Get()
}

func Set(key string, value []byte) (result []byte, err error) {
	return redis.Bytes(db.Get().Do("SET", key, value))
}

func SetEx(key string, value []byte, timeout time.Duration) (result []byte, err error) {
	return redis.Bytes(db.Get().Do("SETEX", key, timeout.Seconds(), value))
}

func Get(key string) (value []byte, err error) {
	return redis.Bytes(db.Get().Do("GET", key))
}

func Del(key string) (err error) {
	_, err = db.Get().Do("DEL", key)
	if err != nil {
		return err
	}

	return
}
