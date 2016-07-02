package game

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"strings"
)

type Console struct {
	Name         string
	Manufacturer string
	Year         int
	User         auth.User
	Has          bool
	HasManual    bool
	HasBox       bool
	Rating       int
	Review       string
	Want         bool
	WantGames    bool
	Picture      string
}

var (
	stmtGetConsoles      *sql.Stmt
	stmtGetUserConsoles  *sql.Stmt
	stmtGetConsole       *sql.Stmt
	stmtUpdateConsole    *sql.Stmt
	stmtUpdateHasConsole *sql.Stmt
	stmtDeleteConsole    *sql.Stmt
	db                   *sql.DB
)

func ConsoleDB(d *sql.DB) {
	var err error
	db = d
	stmtGetConsoles, err = u.Sth(db, "select c.name,IFNULL(c.manufacturer,''),IFNULL(c.year,0),IFNULL(c.picture,'') from consoles as c")
	if err != nil {
		glog.Fatalf("u.Sth(db,select id, IFNULL(name,''),IFNULL(manufacturer,''),IFNULL(year,0) from consoles): %s", err)
	}
	stmtGetConsole, err = u.Sth(db, "select consoles.name,IFNULL(manufacturer,''),IFNULL(year,0),IFNULL(picture,''),IFNULL(user_id,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,''),IFNULL(user_consoles.want,false),IFNULL(user_consoles.wantgames,false) from consoles left join user_consoles on consoles.name=user_consoles.name where consoles.name=? OR (consoles.name=? AND user_consoles.user_id=?)")
	if err != nil {
		glog.Fatalf("u.Sth(db,select id, IFNULL(name,''),IFNULL(manufacturer,''),IFNULL(year,0) from consoles where id=?): %s", err)
	}
	stmtUpdateConsole, err = u.Sth(db, "replace into consoles (name,manufacturer,year,picture) values (?,?,?,?)")
	if err != nil {
		glog.Errorf("u.Sth(db,replace into consoles (name,manufacturer,year) values (?,?,?): %s", err)
	}
	stmtUpdateHasConsole, err = u.Sth(db, "replace into user_consoles (name,user_id,has,manual,box,rating,review,want,wantgames) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		glog.Errorf("u.Sth(db,replace into user_consoles (name,user_id,has,manual,box,rating,review) values (?,?,?,?,?,?,?): %s", err)
	}
	stmtGetUserConsoles, err = u.Sth(db, "select name,IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,''),IFNULL(want,false),IFNULL(wantgames,false) from user_consoles where user_id=?")
	if err != nil {
		glog.Errorf("u.Sth(db,select IFNULL(name,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,'') from user_consoles where user_id=?): %s", err)
	}
	delcon := "delete from consoles where name=? limit 1"
	stmtDeleteConsole, err = u.Sth(db, delcon)
	if err != nil {
		glog.Fatalf("u.Sth(db,%s): %s", delcon, err)
	}

}
func (c Console) String() string {
	return fmt.Sprintf("Name: %s\nManufacturer: %s\nYear: %v\nUserID: %v\nHas: %v\nHasManual: %v\nHasBox: %v\nRating: %v\nReview: %s\nWant: %v\nWantGames: %v\n", c.Name, c.Manufacturer, c.Year, c.User.ID, c.Has, c.HasManual, c.HasBox, c.Rating, c.Review, c.Want, c.WantGames)
}
func (c Console) Save() (e error) {
	if c.Name == "" {
		e = errors.New("Invalid name to Save")
		glog.Errorf("%s", e)
		return e
	}
	_, e = stmtUpdateConsole.Exec(c.Name, c.Manufacturer, c.Year, c.Picture)
	if e != nil {
		glog.Errorf("stmtUpdateConsole.Exec(%s,%s,%v): %s", c.Name, c.Manufacturer, c.Year, e)
		return e
	}
	if c.User.ID > 0 {
		_, e = stmtUpdateHasConsole.Exec(c.Name, c.User.ID, c.Has, c.HasManual, c.HasBox, c.Rating, c.Review, c.Want, c.WantGames)
		if e != nil {
			glog.Errorf("stmtUpdateHasConsole.Exec(%s,%v,%v): %s", c.Name, c.User.ID, e)
			return e
		}
	}
	return e
}
func (c Console) Games() ([]Game, error) {
	gl, err := GetGamesByConsole(c)
	if err != nil {
		glog.Errorf("GetGamesByConsole(%s,%v): %s", c.Name, c.User.ID, err)
	}
	return gl, err
}
func (c Console) OwnedGames() int {
	gl, err := c.Games()
	var og int
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return 0
	}
	for _, g := range gl {
		if g.Has == true {
			og++
		}
	}
	return og
}
func (c Console) UserGames() []Game {
	gl, err := c.Games()
	var rl []Game
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return rl
	}
	for _, g := range gl {
		if g.Has == true {
			rl = append(rl, g)
		}
	}
	return rl
}
func (c Console) Delete() (err error) {
	if !c.User.Admin {
		err = errors.New("game.Console.Delete(): User ! admin")
		return err
	}
	res, err := stmtDeleteConsole.Exec(c.Name)
	if err != nil {
		glog.Errorf("stmtDeleteConsole.exec(%s): %s", c.Name, err)
		return err
	}
	ra, err := res.RowsAffected()
	if ra != 1 {
		e := fmt.Sprintf("game.Consoles.Delete(): Rows affected %v", ra)
		err = errors.New(e)
		return err
	}
	return err
}
func (c Console) TotalGames() int {
	gl, err := c.Games()
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return 0
	}
	return len(gl)
}
func (c Console) StarContent() template.HTML {
	var r string
	for i := 1; i <= 5; i++ {
		var sl string
		if c.Rating >= i {
			sl = "/static/star_on.png"
		} else {
			sl = "/static/star_off.png"
		}
		s := fmt.Sprintf("<img id='star_%s_%v' src='%v' onclick='set_console_rating(\"%v\",\"%s\",%v)'>\n", c.ShortName(), i, sl, c.ShortName(), c.Name, i)
		r = r + " " + s
	}
	return template.HTML(r)
}
func (c Console) ShortName() string {
	return strings.Replace(c.Name, " ", "", -1)
}
func (c Console) WantedGames() (rl []Game, err error) {
	gl, err := c.Games()
	if err != nil {
		glog.Errorf("c.Games(): %s", err)
		return gl, err
	}
	for _, g := range gl {
		if g.Has {
			continue
		}
		if g.Want || c.WantGames {
			rl = append(rl, g)
		}
	}
	return rl, err
}

func GetConsole(name string, user auth.User) (Console, error) {
	var c Console
	var err error
	err = stmtGetConsole.QueryRow(name, name, user.ID).Scan(&c.Name, &c.Manufacturer, &c.Year, &c.Picture, &c.User.ID, &c.Has, &c.HasManual, &c.HasBox, &c.Rating, &c.Review, &c.Want, &c.WantGames)
	c.User = user
	if err != nil {
		glog.Errorf("stmtGetConsoleQueryRow(%s,%v): %s", name, user.ID, err)
	}
	return c, err
}
func GetConsoles(user auth.User) ([]Console, error) {
	var cl []Console
	var cm = make(map[string]Console)
	var um = make(map[string]Console)
	rows, err := stmtGetConsoles.Query()
	if err != nil {
		glog.Errorf("stmtGetConsoles(): %s", err)
		return cl, err
	}
	for rows.Next() {
		var c Console
		rows.Scan(&c.Name, &c.Manufacturer, &c.Year, &c.Picture)
		c.User = user
		cm[c.Name] = c
	}
	if user.ID > 0 {
		rows, err = stmtGetUserConsoles.Query(user.ID)
		if err != nil {
			glog.Errorf("stmtGetUserConsoles(%v): %s", user.ID, err)
			return cl, err
		}
		for rows.Next() {
			var c Console
			rows.Scan(&c.Name, &c.Has, &c.HasManual, &c.HasBox, &c.Rating, &c.Review, &c.Want, &c.WantGames)
			um[c.Name] = c
		}
	}
	for _, c := range cm {
		c.Has = um[c.Name].Has
		c.HasBox = um[c.Name].HasBox
		c.HasManual = um[c.Name].HasManual
		c.Rating = um[c.Name].Rating
		c.Review = um[c.Name].Review
		c.Want = um[c.Name].Want
		c.WantGames = um[c.Name].WantGames
		cl = append(cl, c)
	}
	return cl, err
}

type ConsoleName []Console

func (a ConsoleName) Len() int           { return len(a) }
func (a ConsoleName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ConsoleName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ConsoleMeta struct {
	Console Console
	Games   []Game
}
type ConsoleMetaName []ConsoleMeta

func (a ConsoleMetaName) Len() int           { return len(a) }
func (a ConsoleMetaName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ConsoleMetaName) Less(i, j int) bool { return a[i].Console.Name < a[j].Console.Name }
