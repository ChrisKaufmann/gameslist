package main

import (
	"time"
	"./auth"
	"./game"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
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
	ListEntryHtml	= template.Must(template.ParseFiles("templates/list_entry.html"))
	ConsolesToggle = template.Must(template.ParseFiles("templates/consoles_toggle.html"))
	GamesToggle = template.Must(template.ParseFiles("templates/games_toggle.html"))
	IndentListEntryHtml	= template.Must(template.ParseFiles("templates/indent_list_entry.html"))
	AddConsoleHTML	=	template.Must(template.ParseFiles("templates/add_console.html"))
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
	http.HandleFunc("/main.html", handleMain)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/list/consoles", handleConsoleList)
	http.HandleFunc("/list/games", handleGameList)
	http.HandleFunc("/list/collection", handleMyCollection)
	http.HandleFunc("/toggle/consoles", handleConsolesToggle)
	http.HandleFunc("/toggle/games", handleGamesToggle)
	http.HandleFunc("/collection",handleCollection)
	http.HandleFunc("/console/",handleConsole)
	http.HandleFunc("/thing/", handleThing)
	http.HandleFunc("/mycollection",handleMyCollection)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", handleRoot)
	print("Listening on port " + port + "\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		if err := indexHtml.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/main.html", http.StatusFound)
	}
	t1 := time.Now()
	fmt.Printf("handleRoot %v\n", t1.Sub(t0))
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, _ := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	if err := mainHtml.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t1 := time.Now()
	fmt.Printf("handleMain %v\n", t1.Sub(t0))
}
func handleThing(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	var id string
	u.PathVars(r,"/thing/", &id)
	if id=="" {
		err := errors.New("No id passed to /thing/")
		fmt.Println(err);err.Error();return
	}
	action := r.FormValue("action")
	print("ID:"+id+" action:"+action)
	coll, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	switch action {
		case "toggle":
			t,err :=game.GetThing(id)
			if err != nil {fmt.Println(err);err.Error();return}
			if coll.Have(t) {
				err = coll.Delete(t)
				if err != nil {fmt.Println(err);err.Error();return}
				fmt.Fprintf(w,"white")
			} else {
				coll.Add(t)
				fmt.Fprintf(w,"#aaffa5")
			}
	}
	fmt.Printf("HandleThing %v\n",time.Now().Sub(t0))
}
func handleList(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	t1 := time.Now()
	fmt.Printf("handleList %v\n", t1.Sub(t0))
	return
}
func handleConsolesToggle(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	ConsolesToggle.Execute(w,nil)
	t1 := time.Now()
	fmt.Printf("handleConsolesToggle %v\n", t1.Sub(t0))
	return
}
func handleGamesToggle(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	GamesToggle.Execute(w,nil)
	t1 := time.Now()
	fmt.Printf("handleGamesToggle %v\n", t1.Sub(t0))
	return
}
func handleConsoleList(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	coll, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	var cons []game.MyThing
	switch r.FormValue("filter"){
		case "all":
			tc, err := game.GetAllConsoles()
			if err != nil {fmt.Println(err);err.Error();return}
			cons=coll.MyThingsFromThings(tc)
		case "missing":
			cons, err = coll.MissingConsoles()
			if err != nil {fmt.Println(err);err.Error();return}
		default:
			cons, err = coll.Consoles()
			if err != nil {fmt.Println(err);err.Error();return}
	}
	AddConsoleHTML.Execute(w,nil)
	for _, c := range(cons) {
		p := game.PrintMyThing{c,"white","white","white"}
		if coll.Have(p.Thing){p.Background="#aaffa5"}
		if coll.Have(p.Thing.Box()){p.BoxBackground="#aaffa5"}
		if coll.Have(p.Thing.Manual()){p.ManualBackground="#aaffa5"}
		ListEntryHtml.Execute(w,p)
	}
	t1 := time.Now()
	fmt.Printf("handleConsoleList %v\n", t1.Sub(t0))
}
func handleGameList(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
    loggedin, userID := auth.LoggedIn(w, r)
    if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
    coll, err := game.GetCollection(userID)
	fmt.Printf("handleGameList, after getcollection %v\n", time.Now().Sub(t0))
    if err != nil {fmt.Println(err);err.Error();return}
	var gl []game.MyThing
	switch r.FormValue("filter"){
		case "all":
			tl, err := game.GetAllGames()
			fmt.Printf("handleGameList, after GetAllGames %v\n", time.Now().Sub(t0))
			if err != nil {fmt.Println(err);err.Error();return}
			gl = coll.MyThingsFromThings(tl)
			fmt.Printf("handleGameList, after Mythingsfromthings %v\n", time.Now().Sub(t0))
		case "missing":
			return
		default:
			gl, err = coll.Games()
	}
	fmt.Printf("handleGameList, before execute loop %v\n", time.Now().Sub(t0))
	for _, g := range gl {
		p := game.PrintMyThing{g,"white","white","white"}
		if coll.Have(p.Thing){p.Background="#aaffa5"}
		if coll.Have(p.Thing.Box()){p.BoxBackground="#aaffa5"}
		if coll.Have(p.Thing.Manual()){p.ManualBackground="#aaffa5"}
		ListEntryHtml.Execute(w,p)
	}
	t1 := time.Now()
	fmt.Printf("handleGameList %v\n", t1.Sub(t0))
}
func handleCollection(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleCollection\n")
	//<url>/collection/<id>/<add/remove>/<console/box/manual/note[?note]>
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	_, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	var id		string
	var todo	string
	var thing	string
	u.PathVars(r,"/collection/",&id,&todo,&thing)
	t1 := time.Now()
	fmt.Printf("handleCollection %v\n", t1.Sub(t0))
}
func handleConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	//console/<todo>/<param>
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	_, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	var todo string
	var param string
	u.PathVars(r,"/console/",&todo,&param)
	switch todo {
		case "new":
			_, err := game.AddThing(param,"console")
			if err != nil {
				fmt.Println(err);err.Error();return
			}
		case "newgame":
			console_id := r.FormValue("console_id")
			_, err := game.GetThing(console_id)
			if err != nil {fmt.Println(err);err.Error();return}
	}

	fmt.Printf("handleConsole %v\n", time.Now().Sub(t0))
}
func handleMyCollection(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleMyCollection")
	loggedin, userID := auth.LoggedIn(w, r)
	if ! loggedin { http.Redirect(w,r,"/", http.StatusFound)}
	coll, err := game.GetCollection(userID)
	if err != nil {fmt.Println(err);err.Error();return}
	cons, err := coll.Consoles()
	for _, myc := range(cons) {
		pc := game.PrintMyThing{myc,"white","white","white"}
		if coll.Have(pc.Thing.Box()){pc.BoxBackground="#aaffa5"}
        if coll.Have(pc.Thing.Manual()){pc.ManualBackground="#aaffa5"}
		ListEntryHtml.Execute(w,pc)
		c, err := game.GetThing(myc.ID)
		gl, err := coll.ConsoleGames(c)
		if err != nil {fmt.Println(err);err.Error();return}
		for _, g := range(gl) {
			pg := game.PrintMyThing{g,"#aaffa5","white","white"}
			if coll.Have(pg.Thing.Box()){pg.BoxBackground="#aaffa5"}
			if coll.Have(pg.Thing.Manual()){pg.ManualBackground="#aaffa5"}
			IndentListEntryHtml.Execute(w,pg)
		}
	}
	orphans, err := coll.OrphanGames()
	for _, myg := range orphans {
		ListEntryHtml.Execute(w,myg)
	}
	t1 := time.Now()
	fmt.Printf("handleMyCollection %v\n", t1.Sub(t0))
}
