package main

import (
	"flag"
	"fmt"
	"os"

	"elog.stats/config"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

func loadConf() *config.Config {
	if err := viper.ReadInConfig(); err != nil {
		panic("config: config file not found" + err.Error())
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic("config: unable to load" + err.Error())
	}
	return &conf
}

func main() {
	flag.Usage = usage
	flag.Parse()

	conf := loadConf()

	db := pg.Connect(&pg.Options{
		Addr:     conf.DB.Addr(),
		User:     conf.DB.User,
		Password: conf.DB.Password,
		Database: conf.DB.Database,
	})

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Printf(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
