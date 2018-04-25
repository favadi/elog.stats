package main

import (
	"elog.stats/config"
	"github.com/go-pg/pg"
)

func validDb(db *pg.DB) (err error) {
	_, err = db.Exec("SELECT '1'")
	return
}

func connectDb(conf *config.Config) (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		User:     conf.DB.User,
		Password: conf.DB.Password,
		Database: conf.DB.Database,
		Addr:     conf.DB.Addr(),
	})
	db.PoolStats()
	return db
}

func closeDb(db *pg.DB) {
	if db != nil {
		db.Close()
	}
}
