package main

import (
	"github.com/ChrisKaufmann/goauth"
	"./game"
	"sort"
	"flag"
	"database/sql"
	"strings"
	"errors"
	"fmt"
	"github.com/ChrisKaufmann/easymemcache"
	u "github.com/ChrisKaufmann/goutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/msbranco/goconfig"
	"html/template"
	"net/http"
	"time"
)

var (
	port                 string
	environment          string
	mc                   = easymemcache.New("127.0.0.1:11211")
	cookieName           string
	indexHtml            = template.Must(template.ParseFiles("templates/index-nologin.html"))
	mainHtml             = template.Must(template.ParseFiles("templates/main.html"))
	db                   *sql.DB
	ListEntryHtml        = template.Must(template.ParseFiles("templates/list_entry.html"))
	ConsolesToggle       = template.Must(template.ParseFiles("templates/consoles_toggle.html"))
	GamesToggle          = template.Must(template.ParseFiles("templates/games_toggle.html"))
	IndentListEntryHtml  = template.Must(template.ParseFiles("templates/indent_list_entry.html"))
	ConsoleLinkListEntry = template.Must(template.ParseFiles("templates/console_link_list_entry.html"))
	AddHTML          = template.Must(template.ParseFiles("templates/add.html"))
	ConsoleOnlyEntryHTML = template.Must(template.ParseFiles("templates/console_only_entry.html"))
	TableEntryGameHTML   = template.Must(template.ParseFiles("templates/table_entry_game.html"))
	TableEntryConsoleHTML   = template.Must(template.ParseFiles("templates/table_entry_console.html"))
	SettingsHTML = template.Must(template.ParseFiles("templates/settings.html"))
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
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@"+db_host+"/"+db_name)
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
	game.DB(db)
	game.MemCache(&mc)
	http.HandleFunc("/main.html", handleMain)
	http.HandleFunc("/authorize", auth.HandleAuthorize)
	http.HandleFunc("/settings", handleSettings)
	http.HandleFunc("/oauth2callback", auth.HandleOAuth2Callback)
	http.HandleFunc("/logout", auth.HandleLogout)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/list/games", handleGameList)
	http.HandleFunc("/list/collection", handleMyCollection)
	http.HandleFunc("/toggle/consoles", handleConsolesToggle)
	http.HandleFunc("/toggle/games", handleGamesToggle)
	http.HandleFunc("/collection", handleCollection)
	http.HandleFunc("/console/", handleConsole)
	http.HandleFunc("/thing/", handleThing)
	http.HandleFunc("/mycollection", handleMyCollection)
	http.HandleFunc("/search/", handleSearch)
	http.HandleFunc("/share/", handleShared)
	http.HandleFunc("/demo", handleDemo)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", handleRoot)
	print("Listening on port " + port + "\n")
	http.ListenAndServe("127.0.0.1:"+port, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		fmt.Printf("Not logged in")
		if err := indexHtml.Execute(w, nil); err != nil {
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
	loggedin, _ := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if err := mainHtml.Execute(w, nil); err != nil {
		glog.Errorf("handleMain(): %s", err)
		return
	}
	handleMyCollection(w,r)
	fmt.Fprintf(w,"  </div>	</div>	</body>	</html>")
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
	SettingsHTML.Execute(w, us)
	fmt.Printf("handleSettings %v\n", time.Now().Sub(t0))
}
func handleDemo(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	auth.DemoUser(w,r)
	fmt.Printf("handleDemo %v\n", time.Now().Sub(t0))
}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	fmt.Printf("HandleLogin()")
	var lt string
    u.PathVars(r,"/login/", &lt)
	fmt.Printf("lt: %s", lt)
	err := auth.LoginToken(w,r,lt)
	if err != nil { 
		//http.Redirect(w, r, "/", http.StatusFound)
		glog.Infof("%s", err)
	}
	fmt.Printf("handleLogin %v\n", time.Now().Sub(t0))
}
func handleThing(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var id string
	u.PathVars(r, "/thing/", &id)
	if id == "" {
		glog.Error("No id passed to /thing")
		return
	}
	action := r.FormValue("action")
	print("ID:" + id + " action:" + action)
	coll, err := game.GetCollection(user.ID)
	if err != nil {
		glog.Errorf("handleThing();game.GetCollection(%v): %s", user.ID, err)
		return
	}
	t, err := game.GetThing(id)
	if err != nil {
		glog.Errorf("handleThing:game.GetThing(%s): %s", id, err)
		return
	}
	switch action {
	case "toggle":
		if coll.Have(t) {
			err = coll.Delete(t)
			if err != nil {
				glog.Errorf("handleThing: coll.Have(%s) (trying to toggle unowned)", t.ID)
				return
			}
			fmt.Fprintf(w, "white")
		} else {
			err =coll.Add(t)
			if err!= nil {glog.Errorf("handleThing()coll.Add(%s): %s",t.ID, err);return}
			fmt.Fprintf(w, "#aaffa5")
		}
	case "have":
		err = coll.Add(t)
		if err != nil {glog.Errorf("handleThing()coll.Add(%s): %s", t.ID, err);return}
	case "have_not":
		err = coll.Delete(t)
		if err != nil {
			glog.Errorf("handleThing: coll.Have(%s) (trying to toggle unowned)", t.ID)
			return
		}
	case "setrating":
	    rating := u.Toint(r.FormValue("rating"))
		if rating <1 || rating > 5 { glog.Errorf("Bad rating passed");return}
		tmpl,err := template.New("tmpl").Parse(`{{.}}`)
		pt := coll.GetMyThing(t)
		err =pt.SetRating(rating)
		if err != nil {glog.Errorf("t.SetRating(%s): %s",rating,err) }
		tmpl.Execute(w,coll.GetMyThing(t).StarContent())
	case "get_review_html":
		return
	case "setreview":
		review :=r.FormValue("review")
		tmpl, err := template.New("tmpl").Parse("{{.}}")
		pt := coll.GetMyThing(t)
		err = pt.SetReview(review)
		if err != nil {glog.Errorf("t.SetReview(%s): %s", review, err) }
		tmpl.Execute(w,coll.GetMyThing(t).Review())

	}
	fmt.Printf("HandleThing %v\n", time.Now().Sub(t0))
}
func handleList(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	fmt.Printf("handleList %v\n", time.Now().Sub(t0))
	return
}
func handleConsolesToggle(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	ConsolesToggle.Execute(w, nil)
	fmt.Printf("handleConsolesToggle %v\n", time.Now().Sub(t0))
	return
}
func handleGamesToggle(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	GamesToggle.Execute(w, nil)
	fmt.Printf("handleGamesToggle %v\n", time.Now().Sub(t0))
	return
}
func handleGameList(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	coll, err := game.GetCollection(user.ID)
	fmt.Printf("handleGameList, after getcollection %v\n", time.Now().Sub(t0))
	if err != nil {
		glog.Errorf("handleGameList.game.GetCollection(%v): %s",user.ID,err)
		return
	}
	var gl []game.Thing
	switch r.FormValue("filter") {
	case "all":
		//mytem = ConsoleLinkListEntry
		gl, err = game.GetAllConsoles()
		if err != nil {
			glog.Errorf("handleGameList.game.GetAllConsoles(): %s", err)
			return
		}
		fmt.Fprintf(w,"<table>")
		for _,c := range gl {
			ConsoleLinkListEntry.Execute(w,c)
		}
		fmt.Fprintf(w,"</table>")
		return
	case "console":
		cid := r.FormValue("console_id")
		if cid == "" {
			return
		}
		con, err := game.GetThing(cid)
		if err != nil {
			glog.Errorf("handleGameList.game.GetThing(%s): %s", cid, err)
			return
		}
		gl, err = con.Games()
		if err != nil {
			glog.Errorf("handleGameList.con.Games(): %s", err)
			return
		}
	case "missing":
		return
	default:
		gl, err = coll.Games()
		if err != nil {
			glog.Errorf("handlegameList.coll.Games(): %s", err)
			return
		}
	}
	fmt.Printf("handleGameList, before execute loop %v\n", time.Now().Sub(t0))
	fmt.Fprintf(w, "<table>")
	PrintListOfThings(w,coll,gl)
	fmt.Fprintf(w, "</table>")
	fmt.Printf("handleGameList %v\n", time.Now().Sub(t0))
}
func handleCollection(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleCollection\n")
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	_, err := game.GetCollection(user.ID)
	if err != nil {
		glog.Errorf("handleCollection.game.GetCollection(%v): %s",user.ID, err)
		return
	}
	var id string
	var todo string
	var thing string
	u.PathVars(r, "/collection/", &id, &todo, &thing)
	fmt.Printf("handleCollection %v\n", time.Now().Sub(t0))
}
func handleSearch(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	print("handleCollection\n")
	loggedin,user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	coll, err := game.GetCollection(user.ID)
	if err != nil {
		glog.Errorf("handleSearch.game.GetCollection(%v): %s",user.ID, err)
		return
	}
	ss := r.FormValue("query")
	tl,err := game.Search(ss)
	if err != nil { glog.Errorf("game.Search(%s): %s", ss, err);return }
	fmt.Fprintf(w,"<table>\n")
	PrintListOfThings(w,coll,tl)
	fmt.Fprintf(w,"</table>\n")
	fmt.Printf("handleSearch %v\n", time.Now().Sub(t0))
}
func handleConsole(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	//console/<todo>/<param>
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	_, err := game.GetCollection(user.ID)
	if err != nil {
		glog.Errorf("handleConsole().game.GetCollection(%v): %s", user.ID, err)
		return
	}
	var todo string
	var param string
	u.PathVars(r, "/console/", &todo, &param)
	print("todo: " + todo + ", param: " + param)
	switch todo {
	case "new":
		_, err := game.AddThing(param, "console")
		if err != nil {
			glog.Errorf("handleConsole.game.AddThing(%s, console): %s",param, err)
			return
		}
	case "newgame":
		console_id := r.FormValue("console_id")
		if console_id == "" {
			err = errors.New("No console_id passed to handleConsole(newgame)")
			glog.Error("handleConsole.newgame: No console_id passed")
			return
		}
		game_name := r.FormValue("game_name")
		print("console_id: " + console_id + ", name: " + game_name)
		if game_name == "" {
			glog.Error("No game_name passed to handleConsole(newgame)")
			return
		}
		console, err := game.GetThing(console_id)
		if err != nil {
			glog.Errorf("handleConsole.game.GetThing(%s): %s", console_id, err)
			return
		}
		_, err = console.AddGame(game_name)
		if err != nil {
			glog.Errorf("handleConsole.console.AddGame(%s): %s",game_name, err)
			return
		}
	}
	fmt.Printf("handleConsole %v\n", time.Now().Sub(t0))
}
func handleMyCollection(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	loggedin, user := auth.LoggedIn(w, r)
	if !loggedin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	coll, err := game.GetCollection(user.ID)
	if err != nil {
		glog.Errorf("handleMyCollection.game.GetCollection(%v): %s", user.ID, err)
		return
	}
	cl, err := coll.Things()
	if err != nil { glog.Errorf("handleMyCollection()coll.MyThings(): %s", err);return }
	PrintListOfThings(w,coll,cl)
	fmt.Printf("handleMyCollection %v\n", time.Now().Sub(t0))
}
func PrintListOfThings(w http.ResponseWriter,coll game.Collection,tl []game.Thing) {
	cons, err := game.GetAllConsoles()
	if err != nil {glog.Errorf("PrintListOfThings-game.GetAllConsoles(): %s", err) ;return}
	fmt.Fprintf(w,"<table id='data_table'>")
	curr := "9"
	pttl := coll.GetMyThings(tl)
	sort.Sort(game.ByName(pttl))
	for _, myc := range coll.GetMyThings(cons) {
		TableEntryConsoleHTML.Execute(w,myc)
		for _, t := range pttl {
			//The fc stuff is for printing an anchor 
			fc := strings.ToUpper(t.Name[0:1])
			if fc > curr {
				fmt.Fprintf(w,"<tr><td><a name='"+fc+"' id='"+fc+"'></a></td></tr>\n")
				curr = fc
			}
			if t.ParentID == myc.ID {
				TableEntryGameHTML.Execute(w,t)
			}
		}
	}
	fmt.Fprintf(w,"</table>")
	AddHTML.Execute(w,coll)
}
