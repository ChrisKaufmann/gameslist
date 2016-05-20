package game

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"sort"
)

type Game struct {
	ID          int
	Name        string
	ConsoleName string
	Publisher   string
	Year        int
	User        auth.User
	Has         bool
	HasManual   bool
	HasBox      bool
	Rating      int
	Review      string
}

var (
	stmtGetGame               *sql.Stmt
	stmtUpdateGame            *sql.Stmt
	stmtUpdateUserGame        *sql.Stmt
	stmtGetGameByConsole      *sql.Stmt
	stmtGetUserGamesByConsole *sql.Stmt
	stmtSearchGameNames       *sql.Stmt
	stmtSearchUserGameNames   *sql.Stmt
	stmtInsertGame            *sql.Stmt
	stmtDeleteGame            *sql.Stmt
	stmtGetOwners             *sql.Stmt
	stmtGetUserGames          *sql.Stmt
)

func GameDB(d *sql.DB) {
	var err error
	db = d
	gameselect := "games.id, games.name,games.console_name,IFNULL(publisher,''),IFNULL(year,0),IFNULL(user_id,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,'')"
	sgg := "select " + gameselect + " from games left join user_games on games.id=user_games.game_id where games.id=? OR (games.id=? AND user_games.user_id=?)"
	stmtGetGame, err = u.Sth(db, sgg)
	if err != nil {
		glog.Fatalf("%s: %s", sgg, err)
	}

	sug := "update games set name=?, console_name=?, publisher=?, year=? where id=? limit 1"
	stmtUpdateGame, err = u.Sth(db, sug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sug, err)
	}

	sig := "insert into games (name, console_name, publisher, year) values (?,?,?,?)"
	stmtInsertGame, err = u.Sth(db, sig)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sig, err)
	}

	suug := "replace into user_games (game_id, user_id, has, manual, box, rating, review) values (?,?,?,?,?,?,?)"
	stmtUpdateUserGame, err = u.Sth(db, suug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", suug, err)
	}

	sggbs := "select id, name, console_name, publisher, year from games where console_name=?"
	stmtGetGameByConsole, err = u.Sth(db, sggbs)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sggbs, err)
	}

	sgugbs := "select user_games.game_id, IFNULL(user_id,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,'') from user_games, games where user_games.user_id=? and user_games.game_id=games.id and games.console_name=?"
	stmtGetUserGamesByConsole, err = u.Sth(db, sgugbs)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgugbs, err)
	}

	sgbsn := "select id, name, console_name, publisher, year  from games where name like ?"
	stmtSearchGameNames, err = u.Sth(db, sgbsn)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgbsn, err)
	}

	sgugsearch := "select user_games.game_id, IFNULL(user_id,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,'') from user_games, games where user_games.user_id=? and user_games.game_id=games.id and games.name like ?"
	stmtSearchUserGameNames, err = u.Sth(db, sgugsearch)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgugsearch, err)
	}

	sdelgame := "delete from games where id=? limit 1"
	stmtDeleteGame, err = u.Sth(db, sdelgame)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sdelgame, err)
	}

	sgetowners := "select count(*) from user_games where game_id=?"
	stmtGetOwners, err = u.Sth(db, sgetowners)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgetowners, err)
	}

	sgug := "select games.id, games.name, games.console_name, IFNULL(games.publisher,''), IFNULL(games.year,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,'') from games, user_games where games.id = user_games.game_id and user_games.user_id=?"
	stmtGetUserGames, err = u.Sth(db, sgug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgug, err)
	}
}
func (g Game) String() string {
	return fmt.Sprintf("ID: %v\nName: %s\nConsoleName: %s\nManufacturer: %s\nYear: %v\nUserID: %v\nHas: %v\nHasManual: %v\nHasBox: %v\nRating: %v\nReview: %s\n", g.ID, g.Name, g.ConsoleName, g.Publisher, g.Year, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review)
}
func (g Game) Save() (err error) {
	if g.ID < 1 {
		glog.Errorf("Invalid ID passed to game.Save(): %v", g.ID)
		return errors.New("game.Save: invalid ID")
	}
	if g.User.ID < 1 {
		glog.Errorf("Invalid UserID in game.save(): %v", g.User.ID)
		return errors.New("game.Save: Invalid UserID")
	}
	_, err = stmtUpdateGame.Exec(g.Name, g.ConsoleName, g.Publisher, g.Year, g.ID)
	if err != nil {
		glog.Errorf("stmtUpdateGame.Exec(%v,%s,%s,%s,%v): %s", g.ID, g.Name, g.ConsoleName, g.Publisher, g.Year, err)
		return err
	}
	_, err = stmtUpdateUserGame.Exec(g.ID, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review)
	if err != nil {
		glog.Errorf("stmtUpdateUserGame.Exec(%v,%v,%v,%v,%v,%v,%s): %s", g.ID, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review, err)
		return err
	}
	return err
}
func (g Game) ConsoleSelect() template.HTML {
	var h string
	cl, err := GetConsoles(g.User)
	if err != nil {
		glog.Errorf("GetConsoles(%s): %s", g.User, err)
		return template.HTML("")
	}
	sort.Sort(ConsoleName(cl))
	for _, b := range cl {
		if b.Name == "" {
			continue
		}
		s := ""
		if g.ConsoleName == b.Name {
			s = "selected"
		}
		h = h + fmt.Sprintf(" <option value=\"%s\" %s>%s</option>", b.Name, s, b.Name)
	}
	return template.HTML(h)
}
func (g Game) StarContent() template.HTML {
	var r string
	for i := 1; i <= 5; i++ {
		var sl string
		if g.Rating >= i {
			sl = "/static/star_on.png"
		} else {
			sl = "/static/star_off.png"
		}
		s := fmt.Sprintf("<img id='star_%v_%v' src='%v' onclick='set_game_rating(%v,%v)'>", g.ID, i, sl, g.ID, i)
		r = r + " " + s
	}
	return template.HTML(r)
}
func (g Game) Delete() (err error) {
	if !g.User.Admin {
		err = errors.New("game.Game.Delete(): User != admin")
		return err
	}
	res, err := stmtDeleteGame.Exec(g.ID)
	if err != nil {
		glog.Errorf("stmtDeleteGame.Exec(%v): %s", g.ID, err)
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		glog.Errorf("res.RowsAffected(): %s", err)
	}
	if ra != 1 {
		e := fmt.Sprintf("game.Game.Delete(): Rows affected: %v", ra)
		err = errors.New(e)
		glog.Errorf("%s", err)
	}
	return err
}
func (g Game) Owners() int {
	var retval string
	err := stmtGetOwners.QueryRow(g.ID).Scan(&retval)
	if err != nil {
		glog.Errorf("stmtGetOwners.QueryRow(%v): %s", g.ID, err)
	}
	return u.Toint(retval)
}

func InsertGame(g Game) (Game, error) {
	var err error
	if g.ConsoleName == "" {
		err = errors.New("Game.Insert: Invalid ConsoleName")
		glog.Errorf("%s", err)
		return g, err
	}
	result, err := stmtInsertGame.Exec(g.Name, g.ConsoleName, g.Publisher, g.Year)
	if err != nil {
		glog.Errorf("InsertGame(%s): %s", g, err)
		return g, err
	}
	lid, err := result.LastInsertId()
	if err != nil {
		glog.Errorf("Game.Insert(): %s", err)
		return g, err
	}
	g.ID = int(lid)
	if err := g.Save(); err != nil {
		glog.Errorf("g.Save(): %s", err)
	}
	return g, err
}
func GetGame(id int, user auth.User) (Game, error) {
	var g Game
	var err error
	if user.ID < 1 {
		return g, errors.New("game.GetGame: Invalid UserID")
	}
	err = stmtGetGame.QueryRow(id, id, user.ID).Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year, &g.User.ID, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review)
	g.User = user
	if err != nil {
		e := fmt.Sprintf("GetGame(%v,%s): %s", id, user, err)
		err = errors.New(e)
	}
	return g, err
}
func GetGamesByConsole(con Console) ([]Game, error) {
	var gl []Game
	var err error
	var gm = make(map[int]Game)
	var um = make(map[int]Game)
	rows, err := stmtGetGameByConsole.Query(con.Name)
	if err != nil {
		glog.Errorf("stmtGetGameByConsole(%s): %s", con.Name, err)
		return gl, err
	}
	for rows.Next() {
		var g Game
		rows.Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year)
		g.User = con.User
		gm[g.ID] = g
	}
	rows, err = stmtGetUserGamesByConsole.Query(con.User.ID, con.Name)
	if err != nil {
		glog.Errorf("stmtGetUserGamesByconsole(%v,%s): %s", con.User.ID, con.Name, err)
		return gl, err
	}
	for rows.Next() {
		var g Game
		var dummy int
		rows.Scan(&g.ID, &dummy, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review)
		g.User = con.User
		um[g.ID] = g
	}
	for _, g := range gm {
		g.Has = um[g.ID].Has
		g.HasBox = um[g.ID].HasBox
		g.HasManual = um[g.ID].HasManual
		g.Rating = um[g.ID].Rating
		g.Review = um[g.ID].Review
		g.User = con.User
		gl = append(gl, g)
	}
	return gl, err
}
func SearchGames(ss string, user auth.User) ([]Game, error) {
	var gl []Game
	var err error
	n := "%" + ss + "%"
	var gm = make(map[int]Game)
	var um = make(map[int]Game)
	rows, err := stmtSearchGameNames.Query(n)
	if err != nil {
		glog.Errorf("stmtGetGameByConsole(%s): %s", n, err)
		return gl, err
	}
	for rows.Next() {
		var g Game
		rows.Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year)
		g.User = user
		gm[g.ID] = g
	}
	rows, err = stmtSearchUserGameNames.Query(user.ID, n)
	if err != nil {
		glog.Errorf("stmtGetUserGamesByconsole(%v,%s): %s", user.ID, n, err)
		return gl, err
	}
	for rows.Next() {
		var g Game
		var dummy int
		rows.Scan(&g.ID, &dummy, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review)
		g.User = user
		um[g.ID] = g
	}
	for _, g := range gm {
		g.Has = um[g.ID].Has
		g.HasBox = um[g.ID].HasBox
		g.HasManual = um[g.ID].HasManual
		g.Rating = um[g.ID].Rating
		g.Review = um[g.ID].Review
		g.User = user
		gl = append(gl, g)
	}
	return gl, err
}

type GameName []Game

func (a GameName) Len() int           { return len(a) }
func (a GameName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a GameName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Filter []Game

func (a Filter) Has(tf bool) []Game {
	var gl []Game
	for _, g := range a {
		if g.Has == tf {
			gl = append(gl, g)
		}
	}
	return gl
}
func (a Filter) Box(tf bool) []Game {
	var gl []Game
	for _, g := range a {
		if g.HasBox == tf {
			gl = append(gl, g)
		}
	}
	return gl
}
func (a Filter) Manual(tf bool) []Game {
	var gl []Game
	for _, g := range a {
		if g.HasManual == tf {
			gl = append(gl, g)
		}
	}
	return gl
}

func (a Filter) Request(r *http.Request) []Game {
	if r.FormValue("has") == "true" {
		a = Filter(a).Has(true)
	}
	if r.FormValue("box") == "true" {
		a = Filter(a).Box(true)
	}
	if r.FormValue("manual") == "true" {
		a = Filter(a).Manual(true)
	}
	return a
}

func UserGames(user auth.User) []Game {
	var gl []Game
	var err error
	rows, err := stmtGetUserGames.Query(user.ID)
	if err != nil {
		glog.Errorf("stmtGetUserGames(%v): %s", user.ID, err)
		return gl
	}
	for rows.Next() {
		var g Game
		rows.Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review)
		g.User = user
		gl = append(gl, g)
	}
	return gl

}
