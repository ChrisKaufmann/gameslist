package game

import (
	"database/sql"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/msbranco/goconfig"
	"github.com/stvp/assert"
	"html/template"
	"net/http"
	"net/url"
	"os/exec"
	"testing"
)

func init() {
	initGame()
}
func TestGetGame(t *testing.T) {
	print("GetGame\n")
	seedGame()
	user, err := auth.GetUser(1)
	assert.Nil(t, err, "auth.GetUser(1)")
	g, err := GetGame(1, user)
	assert.Nil(t, err, "GetGame(1,user)")
	assert.Equal(t, 1, g.User.ID, "g.User.ID")
	assert.Equal(t, "game1", g.Name, "g.Name")
	assert.Equal(t, "Nintendo", g.Publisher, "g.Publisher")
	assert.Equal(t, 1990, g.Year, "g.Year")
	assert.Equal(t, true, g.Has, "g.Has")
	assert.Equal(t, true, g.HasBox, "g.HasBox")
	assert.Equal(t, true, g.HasManual, "g.HasManual")
	assert.Equal(t, 3, g.Rating, "g.Rating")
	assert.Equal(t, "is good", g.Review, "g.Review")
	assert.Equal(t, false, g.Want, "Want")
}
func TestGame_Save(t *testing.T) {
	print("Game.Save\n")
	g := gg(t)
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("auth.GetUser(1): %s", err)
	}
	g.Name = "GAME TWO"
	g.ConsoleName = "console2"
	g.Has = false
	g.HasBox = false
	g.HasManual = false
	g.Rating = 5
	g.Review = "is bad"
	g.Publisher = "newman1"
	g.Year = 2000
	g.Want = true
	err = g.Save()
	if err != nil {
		t.Errorf("g.Save(): %s", err)
	}
	d, err := GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(game1,1): %s", err)
	}
	if d.Name != "GAME TWO" {
		t.Errorf("d.Name GAME TWO <=> %s", d.Name)
	}
	if d.Has != false {
		t.Errorf("d.Has false <=> %v", d.Has)
	}
	if d.ConsoleName != "console2" {
		t.Errorf("d.ConsoleName console2 <=> %s", d.ConsoleName)
	}
	if d.HasBox != false {
		t.Errorf("d.HasBox false <=> %v", d.HasBox)
	}
	if d.HasManual != false {
		t.Errorf("d.HasManual false <=> %v", d.HasManual)
	}
	if d.Rating != 5 {
		t.Errorf("d.Rating 5<=>%v", d.Rating)
	}
	if d.Review != "is bad" {
		t.Errorf("d.Review 'is bad' <=> %s", d.Review)
	}
	if d.Publisher != "newman1" {
		t.Errorf("d.Publisher newman1<=>%s", d.Publisher)
	}
	if d.Year != 2000 {
		t.Errorf("d.Year 2000 <=> %v", d.Year)
	}
	assert.Equal(t, true, d.Want, "Want")
}
func TestGame_Delete(t *testing.T) {
	print("Game.Delete\n")
	user := gu(t)
	seedGame()
	g, err := GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(1,user): %s", err)
	}
	print("\tCan't delete if not an admin\n")
	err = g.Delete()
	if err == nil {
		t.Errorf("g.Delete(): %s", err)
	}
	print("\tCan delete if admin\n")
	user.SetAdmin(true)
	user = gu(t)
	g, err = GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(1,user): %s", err)
	}
	err = g.Delete()
	if err != nil {
		t.Errorf("g.Delete(): %s", err)
	}
	g, err = GetGame(1, user)
	if err == nil {
		t.Errorf("Did not delete")
	}
}
func TestGame_Owners(t *testing.T) {
	print("Game.Owners\n")
	g := gg(t)
	if g.Owners() != 1 {
		t.Errorf("g.Owners(): 1 <=> %v", g.Owners())
	}
	g.Has = false
	err := g.Save()
	if err != nil {
		t.Errorf("g.Save(): %s", err)
	}
	if g.Owners() != 0 {
		t.Errorf("g.Owners(): 0 <=> %v", g.Owners())
	}
}
func TestGame_ConsoleSelect(t *testing.T) {
	print("Game.ConsoleSelect\n")
	seedGame()
	user := gu(t)
	g, err := GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(1,user): %s", err)
	}
	css := template.HTML(` <option value="3DS" >3DS</option> <option value="Atari 2600" >Atari 2600</option> <option value="Atari 5200" >Atari 5200</option> <option value="Atari 7800" >Atari 7800</option> <option value="Atari Jaguar" >Atari Jaguar</option> <option value="Atari Lynx" >Atari Lynx</option> <option value="Game Boy" >Game Boy</option> <option value="Game Boy Advance" >Game Boy Advance</option> <option value="Game Boy Color" >Game Boy Color</option> <option value="Game Cube" >Game Cube</option> <option value="Game Gear" >Game Gear</option> <option value="Genesis" >Genesis</option> <option value="Handheld" >Handheld</option> <option value="NES" selected>NES</option> <option value="Nintendo 64" >Nintendo 64</option> <option value="Nintendo DS" >Nintendo DS</option> <option value="PSP" >PSP</option> <option value="Playstation" >Playstation</option> <option value="Playstation 2" >Playstation 2</option> <option value="Playstation 3" >Playstation 3</option> <option value="Playstation 4" >Playstation 4</option> <option value="Playstation Vita" >Playstation Vita</option> <option value="Sega CD" >Sega CD</option> <option value="Sega Dreamcast" >Sega Dreamcast</option> <option value="Sega Saturn" >Sega Saturn</option> <option value="Super Nintendo" >Super Nintendo</option> <option value="Virtual Boy" >Virtual Boy</option> <option value="Wii" >Wii</option> <option value="XBox" >XBox</option> <option value="XBox 360" >XBox 360</option> <option value="XBox One" >XBox One</option>`)
	if g.ConsoleSelect() != css {
		t.Errorf("g.ConsoleSelect()\n'%s'\n<=>'%s'\n", css, g.ConsoleSelect())
	}
}
func TestGame_StarContent(t *testing.T) {
	print("Game.StarContent\n")
	seedGame()
	user := gu(t)
	g, err := GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(1,user): %s", err)
	}
	sc := template.HTML(` <img id='star_1_1' src='/static/star_on.png' onclick='set_game_rating(1,1)'> <img id='star_1_2' src='/static/star_on.png' onclick='set_game_rating(1,2)'> <img id='star_1_3' src='/static/star_on.png' onclick='set_game_rating(1,3)'> <img id='star_1_4' src='/static/star_off.png' onclick='set_game_rating(1,4)'> <img id='star_1_5' src='/static/star_off.png' onclick='set_game_rating(1,5)'>`)
	if g.StarContent() != sc {
		t.Errorf("g.Starcontent()\n'%s'<=>'%s'\n", sc, g.StarContent())
	}
}
func TestGetGamesByConsole(t *testing.T) {
	print("GetGamesByConsole\n")
	c := gsc(t)
	gl, err := GetGamesByConsole(c)
	assert.Nil(t, err, "GetGamesByConsole()")
	assert.Equal(t, 801, len(gl), "GetGamesByConsole")
}

func TestFilter_Box(t *testing.T) {
	print("Filter_Box\n")
	c := gsc(t)
	gl, err := c.Games()
	assert.Nil(t, err, "c.Games()")
	assert.Equal(t, 1, len(Filter(gl).Box(true)), "FilterBox(true)")
	assert.Equal(t, 800, len(Filter(gl).Box(false)), "FilterBox(false)")
}
func TestFilter_Has(t *testing.T) {
	print("Filter_Has\n")
	c := gsc(t)
	gl, err := c.Games()
	assert.Nil(t, err, "c.Games()")
	assert.Equal(t, 1, len(Filter(gl).Has(true)), "FilterHas(true)")
	assert.Equal(t, 800, len(Filter(gl).Has(false)), "FilterHas(false)")
}
func TestFilter_Manual(t *testing.T) {
	print("Filter_Manual\n")
	c := gsc(t)
	gl, err := c.Games()
	assert.Nil(t, err, "c.Games()")
	assert.Equal(t, 1, len(Filter(gl).Manual(true)), "FilterManual(true)")
	assert.Equal(t, 800, len(Filter(gl).Manual(false)), "FilterManual(false)")
}

func TestFilter_Request(t *testing.T) {
	print("Filter_Request\n")
	c := gsc(t)
	gl, err := c.Games()
	assert.Nil(t, err, "c.Games()")
	r := http.Request{}
	uv := url.Values{}
	uv.Add("dummy", "dummy")
	r.Form = uv
	assert.Equal(t, len(Filter(gl).Request(&r)), 801, "FilterRequest()")
	uv.Add("has", "true")
	assert.Equal(t, len(Filter(gl).Request(&r)), 1, "FilterRequest(has=true)")
	uv.Del("has")
	uv.Add("manual", "true")
	assert.Equal(t, len(Filter(gl).Request(&r)), 1, "FilterRequest(manual=true)")
	uv.Del("manual")
	uv.Add("box", "true")
	assert.Equal(t, len(Filter(gl).Request(&r)), 1, "FilterRequest(box=true)")
}
func TestGetGamesByIDS(t *testing.T) {
	print("GetGamesByIDs\n")
	seedGame()
	idl := []int{1, 2, 3, 4, 5}
	user := gu(t)
	gl, err := GetGamesByIDS(idl, user)
	assert.Nil(t, err, "GetGamesByIDS(idl,user)")
	assert.Equal(t, 2, len(gl), "len(gamelist)")
	g := gl[0]
	assert.Equal(t, 1, g.User.ID, "g.User.ID")
	assert.Equal(t, "game1", g.Name, "g.Name")
	assert.Equal(t, "Nintendo", g.Publisher, "g.Publisher")
	assert.Equal(t, 1990, g.Year, "g.Year")
	assert.Equal(t, true, g.Has, "g.Has")
	assert.Equal(t, true, g.HasBox, "g.HasBox")
	assert.Equal(t, true, g.HasManual, "g.HasManual")
	assert.Equal(t, 3, g.Rating, "g.Rating")
	assert.Equal(t, "is good", g.Review, "g.Review")
	assert.Equal(t, false, g.Want, "Want")
}

func TestGetAllWantedGames(t *testing.T) {
	print("GetAllWantedGames\n")
	seedGame()
	user := gu(t)
	c, err := GetConsole("Atari 2600", user)
	assert.Nil(t, err, "GetConsole(Atari 2600,user)")
	print("\tInitial\n")
	wg, err := GetAllWantedGames()
	assert.Nil(t, err, "GetAllWantedGames()")
	assert.Equal(t, 0, len(wg), "GetAllWantedGames")
	print("\tMarking Console WantGames\n")
	c.WantGames = true
	err = c.Save()
	assert.Nil(t, err, "c.Save()")
	wg, err = GetAllWantedGames()
	assert.Nil(t, err, "GetAllWantedGames()")
	assert.Equal(t, 520, len(wg), "GetAllWantedGames()")
	print("\tMarking individual games\n")
	g, err := GetGame(1, user)
	assert.Nil(t, err, "GetGame(1,user)")
	g.Want = true
	err = g.Save()
	assert.Nil(t, err, "g.Save()")
	wg, err = GetAllWantedGames()
	assert.Nil(t, err, "GetAllWantedGames()")
	assert.Equal(t, 521, len(wg), "GetAllWantedGames()")

}

func TestSearchGames(t *testing.T) {
	print("SearchGames\n")
	seedGame()
	user := gu(t)
	var ng Game
	ng.User = user
	ng.Name = "testsearchgame"
	ng.ConsoleName = "NES"
	ng.Has = true
	ng.HasBox = true
	ng.HasManual = true
	ng.Rating = 3
	ng.Publisher = "mypublisher"
	ng.Review = "es muy bueno"
	_, err := InsertGame(ng)
	if err != nil {
		t.Errorf("ng.Save(): %s", err)
	}
	gl, err := SearchGames("testsearchgame", user)
	if err != nil {
		t.Errorf("SearchGames(testsearchgame,user): %s", err)
	}
	if len(gl) != 1 {
		t.Errorf("length SearchGames(testsearchgame,user): 1 <=> %v", len(gl))
	}
	g := gl[0]
	if g.Has != true {
		t.Errorf("g.Has true <=> %v", g.Has)
	}
	if g.HasBox != true {
		t.Errorf("g.HasBox true <=> %v", g.HasBox)
	}
	if g.HasManual != true {
		t.Errorf("g.HasManual true <=> %v", g.HasManual)
	}
	if g.Name != "testsearchgame" {
		t.Errorf("g.Name testsearchgame <=> %v", g.Name)
	}
	if g.ConsoleName != "NES" {
		t.Errorf("g.ConsoleName NES <=> %s", g.ConsoleName)
	}
	if g.User.ID != 1 {
		t.Errorf("g.User.ID 1 <=> %v", g.User.ID)
	}
	if g.Rating != 3 {
		t.Errorf("g.Rating 3 <=> %v", g.Rating)
	}
	if g.Review != "es muy bueno" {
		t.Errorf("g.Review 'es muy bueno' <=> %v", g.Review)
	}
	if g.Publisher != "mypublisher" {
		t.Errorf("g.Publisher 'mypublisher' <=> %v", g.Publisher)
	}
	sl, err := SearchGames("mario", user)
	if err != nil {
		t.Errorf("SearchGames(mario,user): %s", err)
	}
	egl := 9
	if len(sl) != egl {
		t.Errorf("len(SearchGames(mario)) %v <=> %v", egl, len(sl))
	}
}
func TestInsertGame(t *testing.T) {
	print("InsertGame\n")
	seedGame()
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("GetUser(1): %s", err)
	}
	con, err := GetConsole("NES", user)
	if err != nil {
		t.Errorf("GetConsole(NES,user): %s", err)
	}
	nesl, err := GetGamesByConsole(con)
	if err != nil {
		t.Errorf("GetGamesByConsole(): %s", err)
	}
	cl := len(nesl)
	var g Game
	g.Name = "MyNewGame"
	g.ConsoleName = "NES"
	g2, err := InsertGame(g)
	if err != nil {
		t.Errorf("InsertGame(%s): %s", g, err)
	}
	if g2.ConsoleName != "NES" {
		t.Errorf("g2.ConsoleName NES <=> %s", g2.ConsoleName)
	}
	if g2.Name != "MyNewGame" {
		t.Errorf("g2.Name MyNewGame <=> %s", g2.Name)
	}
	if g2.ID < 1 {
		t.Errorf("g2.id: %v", g2.ID)
	}
	cln, err := GetGamesByConsole(con)
	if err != nil {
		t.Errorf("GetGamesByConsole(con): %s", err)
	}
	c2l := len(cln)
	if c2l != cl+1 {
		t.Errorf("New len (%v) != old len(%v)+1", c2l, cl)
	}
}
func initGame() {
	g, err := goconfig.ReadConfigFile("test_config")
	db_name, err := g.GetString("DB", "db")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_host, err := g.GetString("DB", "host")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_user, err := g.GetString("DB", "user")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_pass, err := g.GetString("DB", "pass")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+")/"+db_name)
	if err != nil {
		panic(err)
	}
	out, err := exec.Command("sh", "create_tables.sh").Output()
	if err != nil {
		glog.Errorf("sh create_tables.sh %s\n%s", err, out)
		return
	}

	ConsoleDB(db)
	GameDB(db)
}
func seedGame() {
	out, err := exec.Command("sh", "seed.sh").Output()
	if err != nil {
		glog.Errorf("sh seed.sh %s\n%s", err, out)
		return
	}
}
func gg(t *testing.T) (c Game) {
	seedGame()
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("GetUser(1): %s", err)
	}
	c, err = GetGame(1, user)
	if err != nil {
		glog.Errorf("GetGame(1,1): %s", err)
	}
	return c
}
func dump() {
	out, err := exec.Command("sh", "dump.sh").Output()
	fmt.Printf("%s\n%s\n", out, err)
}
