package main

import (
	"./auth"
	"./game"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"
	"html/template"
	"net/http"
)

var (
	port        string
	environment string
	//mc          = easymemcache.New("127.0.0.1:11211")
	cookieName string
	userName   string
	indexHtml  = template.Must(template.ParseFiles("templates/index-nologin.html"))
	db         *sql.DB
)

func init() {
	var err error
	c, err := goconfig.ReadConfigFile("config")
	if port, err = c.GetString("Web", "port"); err != nil {
		err.Error()
	}
	if environment, err = c.GetString("Web", "environment"); err != nil {
		err.Error()
	}
	db_name, err := c.GetString("DB", "db")
	if err != nil {
		err.Error()
	}
	db_host, err := c.GetString("DB", "host")
	if err != nil {
		err.Error()
	}
	db_user, err := c.GetString("DB", "user")
	if err != nil {
		err.Error()
	}
	db_pass, err := c.GetString("DB", "pass")
	if err != nil {
		err.Error()
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@"+db_host+"/"+db_name)
	if err != nil {
		panic(err)
	}
	cookieName = "gameslist_auth_" + environment
	auth.Environment(environment)
}

func main() {
	defer db.Close()
	auth.DB(db)
	game.DB(db)
	http.HandleFunc("/", handleRoot)
	print("Listening on port " + port + "\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_ = game.GetAllConsoles()
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		if err := indexHtml.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/main", http.StatusFound)
	}
}
