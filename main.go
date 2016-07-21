package main

import (
	"./game"
	"database/sql"
	"flag"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	//	"github.com/ChrisKaufmann/easymemcache"
	u "github.com/ChrisKaufmann/goutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/msbranco/goconfig"
	"html/template"
	"net/http"
	"sort"
	"time"
)

var (
	port        string
	environment string
	cookieName  string
	//	mc                   = easymemcache.New("127.0.0.1:11211")
	db   *sql.DB
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
)

func init() {
	var err error
	flag.Parse()
	c, err := goconfig.ReadConfigFile("config")
	if port, err = c.GetString("Web", "port"); err != nil {
		glog.Fatalf("init(): config.GetString(web.portal): %s", err)
	}
	if environment, err = c.GetString("Web", "environment"); err != nil {
		glog.Fatalf("init(): config.GetString(Web.environment): %s", err)
	}
	db_name, err := c.GetString("DB", "db")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.db): %s", err)
	}
	db_host, err := c.GetString("DB", "host")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.Host): %s", err)
	}
	db_user, err := c.GetString("DB", "user")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.user): %s", err)
	}
	db_pass, err := c.GetString("DB", "pass")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.pass): %s", err)
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+")/"+db_name)
	if err != nil {
		glog.Fatalf("Init():sql.Open(mysql, %s:%s@%s/%s: %s", db_user, db_pass, db_host, db_name, err)
	}
	cookieName = "gameslist_auth_" + environment
	auth.Config("config")
	auth.CookieName(cookieName)
	auth.Environment(environment)
}

func main() {
	defer db.Close()
	auth.DB(db)
	game.GameDB(db)
	game.ConsoleDB(db)
	//	game.MemCache(&mc)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/main.html", handleMain)
	http.HandleFunc("/authorize", auth.HandleAuthorize)
	http.HandleFunc("/settings", handleSettings)
	http.HandleFunc("/oauth2callback", auth.HandleGoogleCallback)
	http.HandleFunc("/logout", auth.HandleLogout)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/search/", handleSearch)
	http.HandleFunc("/share/", handleShared)
	http.HandleFunc("/console/", handleConsole)
	http.HandleFunc("/demo", handleDemo)
	http.HandleFunc("/add/", handleAdd)
	http.HandleFunc("/edit/", handleEdit)
	http.HandleFunc("/set/game/", handleSetGame)
	http.HandleFunc("/set/console/", handleSetConsole)
	http.HandleFunc("/collection", handleCollection)
	http.HandleFunc("/wanted/", handleWanted)
	http.HandleFunc("/", handleRoot)
	print("Listening on port " + port + "\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	fmt.Printf("path: %s\n", r.URL)
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		fmt.Printf("Not logged in")
		if err := tmpl.ExecuteTemplate(w, "index-nologin", nil); err != nil {
			glog.Errorf("HandleRoot(): %s", err)
			return
		}
	} else {
		http.Redirect(w, r, "/main.html", http.StatusFound)
		return
	}
	fmt.Printf("handleRoot %v\n", time.Now().Sub(t0))
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	cl, err := game.GetConsoles(user)
	if err != nil {
		glog.Errorf("game.GetConsoles(%v): %s", user, err)
	}
	sort.Sort(game.ConsoleName(cl))
	type Meta struct {
		User         auth.User
		Consoles     []game.Console
		CheapestGame game.Game
		Search       string
	}
	var m Meta
	m.User = user
	m.Consoles = cl
	wl, err := game.UserWantedGames(user)
	if err != nil {
		glog.Errorf("game.UserWantedGames(user): %s", err)
		return
	}
	m.CheapestGame = game.Filter(wl).Cheapest()
	if err := tmpl.ExecuteTemplate(w, "main_html", m); err != nil {
		glog.Errorf("tmpl.ExecuteTemplate(w, main_html, cl): %s", err)
	}
	fmt.Printf("path: %s\n", r.URL)
	fmt.Printf("handleMain %v\n", time.Now().Sub(t0))
}
func handleSettings(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, us := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Printf("User:  %s", us)
	if err := tmpl.ExecuteTemplate(w, "settings", us); err != nil {
		glog.Errorf("tmpl.ExecuteTemplate(w, settings, user): %s", err)
		return
	}
	fmt.Printf("handleSettings %v\n", time.Now().Sub(t0))
}
func handleDemo(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	auth.DemoUser(w, r)
	fmt.Printf("handleDemo %v\n", time.Now().Sub(t0))
}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	fmt.Printf("HandleLogin()")
	var lt string
	u.PathVars(r, "/login/", &lt)
	fmt.Printf("lt: %s", lt)
	err := auth.LoginToken(w, r, lt)
	if err != nil {
		glog.Infof("%s", err)
	}
	fmt.Printf("handleLogin %v\n", time.Now().Sub(t0))
}
func handleCollection(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	type Meta struct {
		User         auth.User
		Consoles     []game.Console
		CheapestGame game.Game
		Search       string
	}
	cl, err := game.GetConsoles(user)
	if err != nil {
		glog.Errorf("game.GetConsoles(user): %s", err)
		return
	}
	var m Meta
	m.User = user
	wl, err := game.UserWantedGames(user)
	if err != nil {
		glog.Errorf("game.UserWantedGames(user): %s", err)
		return
	}
	m.CheapestGame = game.Filter(wl).Cheapest()

	sort.Sort(game.ConsoleName(cl))
	for _, c := range cl {
		if c.Has || c.OwnedGames() > 0 {
			m.Consoles = append(m.Consoles, c)
		}
	}
	if err := tmpl.ExecuteTemplate(w, "collection", m); err != nil {
		glog.Errorf("ExecuteTemplate(w, collection, m): %s", err)
		return
	}

	fmt.Printf("handleCollection %v\n", time.Now().Sub(t0))
}
func handleWanted(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	type Meta struct {
		User         auth.User
		Consoles     []game.Console
		CheapestGame game.Game
		Url          string
		Has          bool
		Box          bool
		Manual       bool
		Search       string
	}
	var m Meta
	m.User = user
	m.Url = "/wanted/"
	wl, err := game.UserWantedGames(user)
	if err != nil {
		glog.Errorf("game.UserWantedGames(user): %s", err)
		return
	}
	switch r.FormValue("sort") {
	case "soonest":
		sort.Sort(game.GameByEndingSoonest(wl))
	}

	fmt.Printf("handleWanted %v\n", time.Now().Sub(t0))
}
func handleConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	//console/<console name>
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	type Meta struct {
		User         auth.User
		Consoles     []game.Console
		CheapestGame game.Game
		Url          string
		Has          bool
		Box          bool
		Manual       bool
		Search       string
	}
	var m Meta
	m.User = user
	wl, err := game.UserWantedGames(user)
	if err != nil {
		glog.Errorf("game.UserWantedGames(user): %s", err)
		return
	}
	m.CheapestGame = game.Filter(wl).Cheapest()

	if err := tmpl.ExecuteTemplate(w, "main_html", m); err != nil {
		glog.Errorf("Execute main_html: %s", err)
		return
	}
	var cname string
	u.PathVars(r, "/console/", &cname)
	m.Url = "/console/" + cname + "?"

	c, err := game.GetConsole(cname, user)
	fmt.Printf("console=%s", c)
	if err != nil {
		glog.Errorf("game.GetConsole")
	}
	var mycl []game.Console
	mycl = append(mycl, c)
	if err := tmpl.ExecuteTemplate(w, "consoles_list", mycl); err != nil {
		glog.Errorf("ExecuteTemplate(console_entry): %s", err)
		return
	}
	if err != nil {
		glog.Errorf("ConsoleOnlyEntryHTML.Execute(w,myc): %s", err)
	}
	cg, err := c.Games()
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return
	}
	switch r.FormValue("sort") {
	default:
		sort.Sort(game.GameName(cg))
	}
	if r.FormValue("has") == "true" {
		m.Has = true
	}
	if r.FormValue("box") == "true" {
		m.Box = true
	}
	if r.FormValue("manual") == "true" {
		m.Manual = true
	}
	cg = game.Filter(cg).Request(r)
	if err := tmpl.ExecuteTemplate(w, "filter", m); err != nil {
		glog.Errorf("tmpl.ExecuteTemplate(w,filter,meta): %s", err)
	}
	if err := tmpl.ExecuteTemplate(w, "games_list", cg); err != nil {
		glog.Errorf("tmpl.ExecuteTemplate(w, games_list, cg): %s", err)
		return
	}
	fmt.Printf("user: %s %v", user.Email, user)
	fmt.Printf("handleConsole %v\n", time.Now().Sub(t0))
}
func handleAdd(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if !user.Admin {
		glog.Errorf("user %s attempting to add, not an admin", user)
		return
	}
	var cg string
	u.PathVars(r, "/add/", &cg)
	type Meta struct {
		User     auth.User
		Message  template.HTML
		Consoles []game.Console
		Search   string
	}
	var meta Meta
	meta.User = user
	cl, err := game.GetConsoles(user)
	meta.Consoles = cl
	sort.Sort(game.ConsoleName(cl))
	if err != nil {
		glog.Errorf("game.GetConsoles(%s): %s", user, err)
	}
	switch cg {
	case "console":
		var c game.Console
		c.Name = r.PostFormValue("name")
		c.Manufacturer = r.PostFormValue("manufacturer")
		c.Year = u.Toint(r.PostFormValue("year"))
		c.User = user
		c.Picture = r.PostFormValue("picture")
		err := c.Save()
		if err != nil {
			glog.Errorf("Saving console (%s): %s", c, err)
			meta.Message = template.HTML(err.Error())
		} else {
			meta.Message = template.HTML("Saved")
		}
	case "game":
		var g game.Game
		g.Name = r.PostFormValue("name")
		g.Publisher = r.PostFormValue("publisher")
		g.Year = u.Toint(r.PostFormValue("year"))
		g.ConsoleName = r.PostFormValue("console")
		g.User = user
		g, err := game.InsertGame(g)
		if err != nil {
			glog.Errorf("Saving game (%s): %s", g, err)
			meta.Message = "Error saving"
		} else {
			meta.Message = template.HTML("Saved: <a href='/game/" + u.Tostr(g.ID) + "'>Game</a>")
		}
	}
	if err := tmpl.ExecuteTemplate(w, "add", meta); err != nil {
		glog.Errorf("ExecuteTemplate(w,add,meta): %s", err)
	}
	fmt.Printf("handleAdd %v\n", time.Now().Sub(t0))
}
func handleEdit(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if !user.Admin {
		glog.Errorf("user %s attempting to add, not an admin", user)
		return
	}
	var cg string
	u.PathVars(r, "/edit/", &cg)
	switch cg {
	case "console":
		c, err := game.GetConsole(r.FormValue("name"), user)
		if err != nil {
			glog.Errorf("GetConsole(%s,user): %s", r.PostFormValue("name"), err)
		}
		type Meta struct {
			User         auth.User
			Message      template.HTML
			Console      game.Console
			CheapestGame game.Game
			Search       string
		}
		var meta Meta
		meta.User = user
		meta.Console = c
		wl, err := game.UserWantedGames(user)
		if err != nil {
			glog.Errorf("game.UserWantedGames(user): %s", err)
			return
		}
		meta.CheapestGame = game.Filter(wl).Cheapest()

		switch r.PostFormValue("action") {
		case "submit":
			c.Name = r.PostFormValue("name")
			c.Manufacturer = r.PostFormValue("manufacturer")
			c.Year = u.Toint(r.PostFormValue("year"))
			c.Picture = r.PostFormValue("picture")
			if err := c.Save(); err != nil {
				glog.Errorf("Saving console (%s): %s", c, err)
				meta.Message = template.HTML(err.Error())
			} else {
				meta.Message = template.HTML("Saved")
			}
			meta.Console = c
			if err := tmpl.ExecuteTemplate(w, "edit_console", meta); err != nil {
				glog.Errorf("ExecuteTemplate(w,edit_console,meta): %s", err)
			}
		case "delete":
			if err := tmpl.ExecuteTemplate(w, "really_delete_console", meta); err != nil {
				glog.Errorf("ExecuteTemplate(w,really_delete_console,meta): %s", err)
				return
			}
		case "reallydelete":
			err := c.Delete()
			if err != nil {
				glog.Errorf("c.Delete(): %s", err)
			}
			http.Redirect(w, r, "/", http.StatusFound)
		default:
			if err := tmpl.ExecuteTemplate(w, "edit_console", meta); err != nil {
				glog.Errorf("ExecuteTemplate(w,edit_console,meta): %s", err)
			}
		}
	case "game":
		id := u.Toint(r.FormValue("id"))
		g, err := game.GetGame(id, user)
		if err != nil {
			glog.Errorf("game.GetGame(%v,user): %s", id, err)
			return
		}
		type Meta struct {
			User         auth.User
			CheapestGame game.Game
			Game         game.Game
			Message      template.HTML
			Search       string
		}
		var meta Meta
		meta.User = user
		meta.Game = g
		wl, err := game.UserWantedGames(user)
		if err != nil {
			glog.Errorf("game.UserWantedGames(user): %s", err)
			return
		}
		meta.CheapestGame = game.Filter(wl).Cheapest()
		switch r.PostFormValue("action") {
		case "submit":
			g.Name = r.PostFormValue("name")
			g.Publisher = r.PostFormValue("publisher")
			g.Year = u.Toint(r.PostFormValue("year"))
			g.ConsoleName = r.PostFormValue("console")
			if err := g.Save(); err != nil {
				glog.Errorf("Saving game (%s): %s", g, err)
				meta.Message = "Error saving"
			} else {
				fmt.Printf("Saved game: %s", g)
				meta.Message = template.HTML("Saved: <a href='/game/" + u.Tostr(g.ID) + "'>Game</a>")
			}
			meta.Game = g
		case "delete":
			if err := tmpl.ExecuteTemplate(w, "really_delete_game", meta); err != nil {
				glog.Errorf("ExecuteTemplate(w,really_delete_game,meta): %s", err)
				return
			}
		case "reallydelete":
			err := g.Delete()
			if err != nil {
				glog.Errorf("g.Delete(): %s", err)
			}
			http.Redirect(w, r, "/", http.StatusFound)

		default:
			meta.Message = template.HTML("in default")
		}
		if err := tmpl.ExecuteTemplate(w, "edit_game", meta); err != nil {
			glog.Errorf("ExecuteTemplate(w,edit_game, meta): %s", err)
			return
		}
	}
	fmt.Printf("handleEdit %v\n", time.Now().Sub(t0))
}

func handleSetConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	name := r.FormValue("name")
	if len(name) < 1 {
		glog.Errorf("Invalid Name passed to set/console: %s", name)
		return
	}
	if r.FormValue("action") == "" {
		glog.Errorf("No action passed to set/console/")
		return
	}
	c, err := game.GetConsole(name, user)
	if err != nil {
		glog.Errorf("GetConsole(%v,%v): %s", name, user, err)
	}
	fmt.Printf("in HandleSetConsole: Console:\n%s--------\n", c)
	switch r.FormValue("action") {
	case "setrating":
		c.Rating = u.Toint(r.FormValue("rating"))
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
		fmt.Fprintf(w, "%s", c.StarContent())
		return
	case "have":
		c.Has = true
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "have_not":
		c.Has = false
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "has_manual":
		c.HasManual = true
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "hasnot_manual":
		c.HasManual = false
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "has_box":
		c.HasBox = true
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "hasnot_box":
		c.HasBox = false
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "want":
		c.Want = true
		err := c.Save()
		if err != nil {
			glog.Error("g.Save(): %s", err)
		}
	case "wantnot":
		c.Want = false
		err := c.Save()
		if err != nil {
			glog.Error("g.Save(): %s", err)
		}
	case "wantgames":
		c.WantGames = true
		err := c.Save()
		if err != nil {
			glog.Error("g.Save(): %s", err)
		}
	case "wantnotgames":
		c.WantGames = false
		err := c.Save()
		if err != nil {
			glog.Error("g.Save(): %s", err)
		}
	case "set_review":
		c.Review = r.FormValue("review")
		err := c.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
		fmt.Fprintf(w, "%s", c.Review)
	default:
		glog.Errorf("Invalid action passed to set console: %s", r.FormValue("action"))
	}
	fmt.Printf("handleSetConsole %v\n", time.Now().Sub(t0))
}
func handleSetGame(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var id int
	if r.FormValue("id") == "" {
		glog.Errorf("No ID passed to SetGame")
		return
	}
	if r.FormValue("action") == "" {
		glog.Errorf("No action passed to SetGame")
		return
	}
	id = u.Toint(r.FormValue("id"))
	g, err := game.GetGame(id, user)
	if err != nil {
		glog.Errorf("GetGame(%v,%v): %s", id, user, err)
	}
	switch r.FormValue("action") {
	case "setrating":
		g.Rating = u.Toint(r.FormValue("rating"))
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
		fmt.Fprintf(w, "%s", g.StarContent())
	case "have":
		g.Has = true
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "have_not":
		g.Has = false
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "want":
		g.Want = true
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "wantnot":
		g.Want = false
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "has_manual":
		g.HasManual = true
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "hasnot_manual":
		g.HasManual = false
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "has_box":
		g.HasBox = true
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "hasnot_box":
		g.HasBox = false
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
	case "set_review":
		g.Review = r.FormValue("review")
		err := g.Save()
		if err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
		fmt.Fprintf(w, "%s", g.Review)
	default:
		glog.Errorf("Invalid action passed to set game: %s", r.FormValue("action"))
	}
	fmt.Printf("handleSetGame %v\n", time.Now().Sub(t0))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleCollection\n")
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	cml := make(map[string][]game.Game)
	type Meta struct {
		ConsoleMeta  []game.ConsoleMeta
		CheapestGame game.Game
		User         auth.User
		Search       string
		Url          string
		Has          bool
		Box          bool
		Manual       bool
	}
	var meta Meta
	meta.User = user
	wl, err := game.UserWantedGames(user)
	if err != nil {
		glog.Errorf("game.UserWantedGames(user): %s", err)
		return
	}
	meta.CheapestGame = game.Filter(wl).Cheapest()
	meta.Url = "/search/?query=" + r.FormValue("query") + "&"
	gl, err := game.SearchGames(r.FormValue("query"), user)
	if err != nil {
		glog.Errorf("game.SearchGames(%s, user): %s", r.FormValue("query"), err)
	}
	sort.Sort(game.GameName(gl))
	if r.FormValue("has") == "true" {
		meta.Has = true
	}
	if r.FormValue("box") == "true" {
		meta.Box = true
	}
	if r.FormValue("manual") == "true" {
		meta.Manual = true
	}
	gl = game.Filter(gl).Request(r)
	for _, g := range gl {
		cml[g.ConsoleName] = append(cml[g.ConsoleName], g)
	}
	cl, err := game.GetConsoles(user)
	var sm []game.ConsoleMeta
	for _, c := range cl {
		if len(cml[c.Name]) > 0 {
			var cm game.ConsoleMeta
			cm.Console = c
			cm.Games = cml[c.Name]
			sm = append(sm, cm)
		}
	}
	sort.Sort(game.ConsoleMetaName(sm))
	meta.ConsoleMeta = sm
	meta.Search = r.FormValue("query")
	if err := tmpl.ExecuteTemplate(w, "search", meta); err != nil {
		glog.Errorf("ExecuteTemplate: %s", err)
	}
	fmt.Printf("handleSearch %v\n", time.Now().Sub(t0))
}
