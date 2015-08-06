package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"
)

var (
	db      *sql.DB
	db_name string
	db_host string
	db_user string
	db_pass string
)

func init() {
	var err error
	c, err := goconfig.ReadConfigFile("config")
	if err != nil {
		err.Error()
	}
	db_name, err = c.GetString("DB", "db")
	if err != nil {
		err.Error()
	}
	db_host, err = c.GetString("DB", "host")
	if err != nil {
		err.Error()
	}
	db_user, err = c.GetString("DB", "user")
	if err != nil {
		err.Error()
	}
	db_pass, err = c.GetString("DB", "pass")
	if err != nil {
		err.Error()
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@"+db_host+"/"+db_name)
	if err != nil {
		panic(err)
	}
}
func sth(s string) *sql.Stmt {
	print(s)
	a, err := db.Prepare(s)
	if err != nil {
		print(s)
		err.Error()
	}
	print("Prepared")
	return a
}
