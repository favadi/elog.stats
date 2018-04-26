package main

import (
	"github.com/go-pg/pg"
	"github.com/txchuyen/elog.stats/config"
)

// TODO: this function is unused, what is it purpose?
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
		// TODO: db.Close() returns an error, handle it
		db.Close()
	}
}
