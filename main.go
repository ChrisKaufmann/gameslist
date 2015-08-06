package main

import (
	"github.com/ChrisKaufmann/easymemcache"
	"github.com/msbranco/goconfig"
	"html/template"
	"net/http"
)

var (
	port        string
	environment string
	mc          = easymemcache.New("127.0.0.1:11211")
	cookieName  string
	userName    string
	indexHtml         = template.Must(template.ParseFiles("templates/index-nologin.html"))
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
	cookieName = "gameslist_auth_" + environment
}

func main() {
	defer db.Close()
	http.HandleFunc("/", handleRoot)
	print("Listening on port "+port+"\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if !loggedIn(w, r) {
		if err := indexHtml.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/main", http.StatusFound)
	}
}
