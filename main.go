package main

import (
	"./auth"
	"./game"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"
	"html/template"
	"net/http"
	u "github.com/ChrisKaufmann/goutils"
	"fmt"
)

var (
	port        string
	environment string
	//mc		= easymemcache.New("127.0.0.1:11211")
	cookieName	string
	indexHtml	= template.Must(template.ParseFiles("templates/index-nologin.html"))
	mainHtml	= template.Must(template.ParseFiles("templates/main.html"))
	db         *sql.DB
	consoleListEntryHtml	= template.Must(template.ParseFiles("templates/consolelist_entry.html"))
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
	http.HandleFunc("/main.html", handleMain)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/list/consoles", handleConsoleList)
	http.HandleFunc("/collection",handleCollection)
	http.HandleFunc("/mycollection",handleMyCollection)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	print("Listening on port " + port + "\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		if err := indexHtml.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/main.html", http.StatusFound)
	}
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	loggedin, _ := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	if err := mainHtml.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func handleList(w http.ResponseWriter, r *http.Request) {
	return
}
func handleConsoleList(w http.ResponseWriter, r *http.Request) {
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	coll, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	cons, err := coll.Consoles()
	for _, c := range(cons) {
		consoleListEntryHtml.Execute(w,c)
	}

}
func handleCollection(w http.ResponseWriter, r *http.Request) {
	//<url>/collection/<id>/<add/remove>/<console/box/manual/note[?note]>
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	_, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	var id		string
	var todo	string
	var thing	string
	u.PathVars(r,"/collection/",&id,&todo,&thing)
}
func handleMyCollection(w http.ResponseWriter, r *http.Request) {
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	coll, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	cons, err := coll.Consoles()
	for _, c := range(cons) {
		consoleListEntryHtml.Execute(w,c)
		gl, err := coll.ConsoleGames(c)
		if err != nil {fmt.Println(err);err.Error();return}
		cons_box, err := c.Box()
		if coll.Have(cons_box) { fmt.Fprintf(w,"\tBox:"+cons_box.Name+"\n") }
		cons_m, err := c.Manual()
		if coll.Have(cons_m) { fmt.Fprintf(w,"\tManual:"+cons_m.Name+"\n") }
		for _, g := range(gl) {
			fmt.Fprintf(w,"\tGame:"+g.Name+"\n")
			g_box, err := g.Box()
			if err != nil {fmt.Println(err);err.Error();return}
			if coll.Have(g_box) {fmt.Fprintf(w, "\t\tBox:"+g_box.Name+"\n") }
			g_man, err := g.Manual()
			if coll.Have(g_man) {fmt.Fprintf(w, "\t\tManual:"+g_man.Name+"\n") }
			if err != nil {fmt.Println(err);err.Error();return}
		}
	}
}
