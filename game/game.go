package game

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ChrisKaufmann/ebay-go"
	"github.com/ChrisKaufmann/goauth"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Game struct {
	ID          int
	Name        string
	ConsoleName string
	Publisher   string
	Year        int
	EbayPrice   float64
	EbayUpdated string
	EbayEnds    string
	EbayURL     string
	User        auth.User
	Has         bool
	HasManual   bool
	HasBox      bool
	Rating      int
	Review      string
	Want        bool
}

var (
	stmtGetGame                 *sql.Stmt
	stmtUpdateGame              *sql.Stmt
	stmtUpdateUserGame          *sql.Stmt
	stmtGetGameByConsole        *sql.Stmt
	stmtGetUserGamesByConsole   *sql.Stmt
	stmtSearchGameNames         *sql.Stmt
	stmtSearchUserGameNames     *sql.Stmt
	stmtInsertGame              *sql.Stmt
	stmtDeleteGame              *sql.Stmt
	stmtGetOwners               *sql.Stmt
	stmtGetUserGames            *sql.Stmt
	stmtGetWantedGamesByConsole *sql.Stmt
	stmtGetWantedGamesByGame    *sql.Stmt
	stmtUserWantedGames         *sql.Stmt
)

func GameDB(d *sql.DB) {
	var err error
	db = d
	gameselect := "games.id, games.name,games.console_name,IFNULL(publisher,''),IFNULL(year,0),IF(ebay_ends > NOW(),ebay_price,0.0),IFNULL(ebay_updated,''),IFNULL(ebay_ends,''),IFNULL(ebay_url,''),IFNULL(user_id,0),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,''),IFNULL(want,false)"
	sgg := "select " + gameselect + " from games left join user_games on games.id=user_games.game_id where games.id=? OR (games.id=? AND user_games.user_id=?)"
	stmtGetGame, err = u.Sth(db, sgg)
	if err != nil {
		glog.Fatalf("%s: %s", sgg, err)
	}

	sug := "update games set name=?, console_name=?, publisher=?, year=?, ebay_price=?, ebay_ends=?, ebay_url=?,ebay_updated=? where id=? limit 1"
	stmtUpdateGame, err = u.Sth(db, sug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sug, err)
	}

	sig := "insert into games (name, console_name, publisher, year) values (?,?,?,?)"
	stmtInsertGame, err = u.Sth(db, sig)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sig, err)
	}

	suug := "replace into user_games (game_id, user_id, has, manual, box, rating, review,want) values (?,?,?,?,?,?,?,?)"
	stmtUpdateUserGame, err = u.Sth(db, suug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", suug, err)
	}

	sggbs := "select id from games where console_name=?"
	stmtGetGameByConsole, err = u.Sth(db, sggbs)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sggbs, err)
	}

	sgbsn := "select id  from games where name like ?"
	stmtSearchGameNames, err = u.Sth(db, sgbsn)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgbsn, err)
	}

	sdelgame := "delete from games where id=? limit 1"
	stmtDeleteGame, err = u.Sth(db, sdelgame)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sdelgame, err)
	}

	sgetowners := "select count(*) from user_games where game_id=? and has=true"
	stmtGetOwners, err = u.Sth(db, sgetowners)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgetowners, err)
	}

	sgug := "select game_id from user_games where user_id=?"
	stmtGetUserGames, err = u.Sth(db, sgug)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sgug, err)
	}

	sql := "select distinct(g.id) from games as g, consoles as c, user_consoles as uc where (uc.wantgames=1 and uc.name=g.console_name)"
	stmtGetWantedGamesByConsole, err = u.Sth(db, sql)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sql, err)
	}

	sql = "select distinct(g.game_id) from user_games as g where want=1"
	stmtGetWantedGamesByGame, err = u.Sth(db, sql)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sql, err)
	}

	sql = "select distinct(g.id) from games as g, user_consoles as uc, user_games as ug where (uc.wantgames=1 and uc.name=g.console_name and uc.user_id=?) or (ug.user_id=? and ug.want=1 and ug.game_id=g.id)"
	stmtUserWantedGames, err = u.Sth(db, sql)
	if err != nil {
		glog.Fatalf("u.Sth(db, %s): %s", sql, err)
	}

}
func (g Game) String() string {
	return fmt.Sprintf("ID: %v\nName: %s\nConsoleName: %s\nManufacturer: %s\nYear: %v\nUserID: %v\nHas: %v\nHasManual: %v\nHasBox: %v\nRating: %v\nReview: %s\nWant: %v\n", g.ID, g.Name, g.ConsoleName, g.Publisher, g.Year, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review, g.Want)
}
func (g Game) Price() string {
	return fmt.Sprintf("%.2f", g.EbayPrice)
}
func (g *Game) Save() (err error) {
	if g.ID < 1 {
		glog.Errorf("Invalid ID passed to game.Save(): %v", g.ID)
		return errors.New("game.Save: invalid ID")
	}
	_, err = stmtUpdateGame.Exec(g.Name, g.ConsoleName, g.Publisher, g.Year, g.EbayPrice, g.EbayEnds, g.EbayURL, g.EbayUpdated, g.ID)
	if err != nil {
		glog.Errorf("stmtUpdateGame.Exec(%v,%s,%s,%s,%v): %s", g.ID, g.Name, g.ConsoleName, g.Publisher, g.Year, err)
		return err
	}
	if g.User.ID > 0 {
		_, err = stmtUpdateUserGame.Exec(g.ID, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review, g.Want)
		if err != nil {
			glog.Errorf("stmtUpdateUserGame.Exec(%v,%v,%v,%v,%v,%v,%s): %s", g.ID, g.User.ID, g.Has, g.HasManual, g.HasBox, g.Rating, g.Review, err)
			return err
		}
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
func (g *Game) UpdateEbay(eb *ebay.EBay) error {
	searchstring := fmt.Sprintf("%s %s", g.ConsoleName, g.Name)
	il, err := eb.Search(searchstring)
	if err != nil {
		glog.Errorf("eb.Search(%s): %s", searchstring, err)
		return err
	}
	l := ebay.LowestPricePlusShipping(il)
	g.EbayPrice = l.Price + l.ShippingPrice
	g.EbayEnds = l.EndTime
	g.EbayURL = l.Url
	fmt.Printf("%s\t%v\t%v\n", g.Name, l.Price, l.ShippingPrice)
	err = g.Save()
	if err != nil {
		glog.Errorf("g.Save(): %s", err)
	}
	return err
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
	if g.User.ID > 0 {
		if err := g.Save(); err != nil {
			glog.Errorf("g.Save(): %s", err)
		}
		return g, err
	}
	return g, err
}
func GetGame(id int, user auth.User) (Game, error) {
	var g Game
	var err error
	err = stmtGetGame.QueryRow(id, id, user.ID).Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year, &g.EbayPrice, &g.EbayUpdated, &g.EbayEnds, &g.EbayURL, &g.User.ID, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review, &g.Want)
	g.User = user
	if err != nil {
		e := fmt.Sprintf("GetGame(%v,%s): %s", id, user, err)
		err = errors.New(e)
	}

	return g, err
}
func GetGamesByConsole(c Console) (gl []Game, err error) {
	var idl []int
	rows, err := stmtGetGameByConsole.Query(c.Name)
	if err != nil {
		glog.Errorf("stmtGetGameByConsole(%s): %s", c.Name, err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		idl = append(idl, id)
	}
	gl, err = GetGamesByIDS(idl, c.User)
	if err != nil {
		glog.Errorf("GetGamesByIDS(id_list, user): %s", err)
	}
	return gl, err
}
func SearchGames(ss string, user auth.User) (gl []Game, err error) {
	var idl []int
	n := "%" + ss + "%"
	rows, err := stmtSearchGameNames.Query(n)
	if err != nil {
		glog.Errorf("stmtGetGameByConsole(%s): %s", n, err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		idl = append(idl, id)
	}
	gl, err = GetGamesByIDS(idl, user)
	if err != nil {
		glog.Errorf("GetGamesByIDS(idl, user): %s", err)
	}
	return gl, err
}
func GetGamesByIDS(id_list []int, user auth.User) (gl []Game, err error) {
	if len(id_list) == 1 {
		g, err := GetGame(id_list[0], user)
		gl = append(gl, g)
		return gl, err
	}
	if len(id_list) == 0 {
		return gl, err
	}
	var id_list_string []string
	for _, id := range id_list {
		id_list_string = append(id_list_string, strconv.Itoa(id))
	}
	ids := strings.Join(id_list_string, ",")
	sql := fmt.Sprintf("select games.id, games.name,games.console_name,IFNULL(publisher,''),IFNULL(year,0),IF(ebay_ends > NOW(),ebay_price,0.0),IFNULL(ebay_updated,''),IFNULL(ebay_ends,''),IFNULL(ebay_url,''),IFNULL(has,false),IFNULL(manual,false),IFNULL(box,false),IFNULL(rating,0),IFNULL(review,''),IFNULL(want,false) from  games left join user_games on games.id=user_games.game_id where games.id in (%s) OR (games.id in (%s) AND user_games.user_id=%v) ;", ids, ids, user.ID)
	if user.ID < 1 {
		sql = fmt.Sprintf("select games.id, games.name,games.console_name,IFNULL(publisher,''),IFNULL(year,0),IFNULL(ebay_price,0.0),IFNULL(ebay_updated,''),IFNULL(ebay_ends,''),IFNULL(ebay_url,''),false,false,false,0,'',false from games where id in (%s)", ids)
	}
	sth, err := u.Sth(db, sql)
	if err != nil {
		glog.Errorf("u.Sth(db, %s): %s", sql, err)
	}
	rows, err := sth.Query()
	if err != nil {
		glog.Errorf("sth.Query(): %s", err)
	}
	for rows.Next() {
		var g Game
		rows.Scan(&g.ID, &g.Name, &g.ConsoleName, &g.Publisher, &g.Year, &g.EbayPrice, &g.EbayUpdated, &g.EbayEnds, &g.EbayURL, &g.Has, &g.HasManual, &g.HasBox, &g.Rating, &g.Review, &g.Want)
		g.User = user
		gl = append(gl, g)
	}
	return gl, err
}

func GetAllWantedGames() (gl []Game, err error) {
	var idl []int
	id_map := make(map[int]int)
	rows, err := stmtGetWantedGamesByConsole.Query()
	if err != nil {
		glog.Errorf("stmtGetWantedGamesByConsole(): %s", err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		id_map[id] = 1
	}
	rows, err = stmtGetWantedGamesByGame.Query()
	if err != nil {
		glog.Errorf("stmtGetWantedGamesByGame(): %s", err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		id_map[id] = 1
	}
	for id, _ := range id_map {
		idl = append(idl, id)
	}
	var user auth.User
	gl, err = GetGamesByIDS(idl, user)
	if err != nil {
		glog.Errorf("GetGamesByIDS(idl,user): %s", err)
	}
	return gl, err
}
func UserWantedGames(user auth.User) (gl []Game, err error) {
	var idl []int
	rows, err := stmtUserWantedGames.Query(user.ID, user.ID)
	if err != nil {
		glog.Errorf("stmtUserWantedGames.Query(%v,%v): %s", user.ID, user.ID, err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		idl = append(idl, id)
	}
	gl, err = GetGamesByIDS(idl, user)
	if err != nil {
		glog.Errorf("GetGamesByIDS(list..., user): %s", err)
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
func (a Filter) Cheapest() Game {
	sort.Sort(GameByPrice(a))
	return a[0]
}

type GameByPrice []Game

func (a GameByPrice) Len() int      { return len(a) }
func (a GameByPrice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a GameByPrice) Less(i, j int) bool {
	if a[i].EbayPrice < 0.01 {
		return false
	}
	if a[j].EbayPrice < 0.01 {
		return true
	}
	return a[i].EbayPrice < a[j].EbayPrice
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
