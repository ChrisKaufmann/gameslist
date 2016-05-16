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
	http.HandleFunc("/oauth2callback", auth.HandleOAuth2Callback)
	http.HandleFunc("/logout", auth.HandleLogout)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/search/", handleSearch)
	http.HandleFunc("/share/", handleShared)
	http.HandleFunc("/console/", handleConsole)
	http.HandleFunc("/demo", handleDemo)
	http.HandleFunc("/set/game/", handleSetGame)
	http.HandleFunc("/set/console/", handleSetConsole)
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
	if err := tmpl.ExecuteTemplate(w, "main_html", cl); err != nil {
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
func handleConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	//console/<console name>
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "main_html", nil); err != nil {
		glog.Errorf("Execute main_html: %s", err)
		return
	}
	var cname string
	u.PathVars(r, "/console/", &cname)
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
	switch r.FormValue("sort") {
	default:
		sort.Sort(game.GameName(cg))
	}
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "games_list", cg); err != nil {
		glog.Errorf("tmpl.ExecuteTemplate(w, games_list, cg): %s", err)
		return
	}
	fmt.Printf("user: %s %v", user.Email, user)
	fmt.Printf("handleConsole %v\n", time.Now().Sub(t0))
}
func handleSetGame(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var id int
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
func handleSetConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	name := r.FormValue("name")
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
	var cl []game.Console
	cl = append(cl, c)
	if err := tmpl.ExecuteTemplate(w, "main_html", cl); err != nil {
		glog.Errorf("ExecuteTemplate(w, consoles_list, cl): %s", err)
	}
	fmt.Printf("handleSetConsole %v\n", time.Now().Sub(t0))
}
func handleSearch(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleCollection\n")
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	/*
		coll, err := game.GetCollection(user)
		if err != nil {
			glog.Errorf("handleSearch.game.GetCollection(%v): %s",user, err)
			return
		}
		ss := r.FormValue("query")
		tl,err := game.Search(ss)
		if err != nil { glog.Errorf("game.Search(%s): %s", ss, err);return }
		fmt.Fprintf(w,"<table>\n")
		//PrintListOfThings(w,coll,tl)
	*/
	cml := make(map[string][]game.Game)
	type ConsoleMeta struct {
		Console game.Console
		Games   []game.Game
	}
	gl, err := game.SearchGames(r.FormValue("query"), user)
	if err != nil {
		glog.Errorf("game.SearchGames(%s, user): %s", r.FormValue("query"), err)
	}
	sort.Sort(game.GameName(gl))
	for _, g := range gl {
		cml[g.ConsoleName] = append(cml[g.ConsoleName], g)
	}
	cl, err := game.GetConsoles(user)
	var sm []ConsoleMeta
	for _, c := range cl {
		if len(cml[c.Name]) > 0 {
			var cm ConsoleMeta
			cm.Console = c
			cm.Games = cml[c.Name]
			sm = append(sm, cm)
		}
	}
	if err := tmpl.ExecuteTemplate(w, "search", sm); err != nil {
		glog.Errorf("ExecuteTemplate: %s", err)
	}
	fmt.Printf("handleSearch %v\n", time.Now().Sub(t0))
}
