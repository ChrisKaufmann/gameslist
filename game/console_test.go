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
	"os/exec"
	"testing"
)

func init() {
	initConsole()
}
func TestGetConsole(t *testing.T) {
	print("GetConsole\n")
	seedConsole()
	user := gu(t)
	c, err := GetConsole("NES", user)
	if err != nil {
		t.Errorf("GetConsole(NES,1): %s", err)
	}
	assert.Equal(t, 1, c.User.ID, "c.UserID")
	assert.True(t, c.Has, "c.Has")
	assert.True(t, c.HasBox, "c.HasBox")
	assert.True(t, c.HasManual, "c.HasManual")
	assert.Equal(t, 3, c.Rating, "c.Rating")
	assert.Equal(t, "is good", c.Review, "c.Review")
	assert.Equal(t, "Nintendo", c.Manufacturer, "c.Manufacturer")
	assert.False(t, c.Want, "c.Want")
	assert.False(t, c.WantGames, "c.WantGames")
}
func TestConsole_Save(t *testing.T) {
	print("Console.Save\n")
	seedConsole()
	c := gsc(t)
	c.Has = false
	c.HasBox = false
	c.HasManual = false
	c.Rating = 5
	c.Review = "is bad"
	c.Manufacturer = "newman1"
	c.Year = 2000
	c.Want = true
	c.WantGames = true
	err := c.Save()
	if err != nil {
		t.Errorf("c.Save(): %s", err)
	}
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("GetUser(1): %s", err)
	}
	d, err := GetConsole("NES", user)
	if err != nil {
		t.Errorf("GetConsole(NES,user): %s", err)
	}
	assert.False(t, d.Has, "console.Has")
	assert.False(t, d.HasBox, "console.HasBox")
	assert.False(t, d.HasManual, "console.HasManual")
	assert.Equal(t, 5, d.Rating, "console.Rating")
	assert.Equal(t, "is bad", d.Review, "console.Review")
	assert.Equal(t, "newman1", d.Manufacturer, "console.Manufacturer")
	assert.Equal(t, 2000, d.Year, "console.Year")
	assert.Equal(t, true, d.Want, "console.Want")
	assert.Equal(t, true, d.WantGames, "console.Wantgames")
}
func TestConsole_Delete(t *testing.T) {
	print("Console.Delete\n")
	seedConsole()
	user := gu(t)
	cl, err := GetConsoles(user)
	if err != nil {
		t.Errorf("GetConsoles(1): %s", err)
	}
	lc := len(cl)
	c, err := GetConsole("NES", user)
	if err != nil {
		glog.Errorf("GetConsole(NES,1): %s", err)
	}
	print("\tCan't delete if not an admin\n")
	if err := c.Delete(); err == nil {
		t.Errorf("c.Delete(): %s", err)
	}
	user.SetAdmin(true)
	user = gu(t)
	c, err = GetConsole("NES", user)
	if err != nil {
		glog.Errorf("GetConsole(NES,1): %s", err)
	}
	print("\tCan delete if an admin\n")
	if err := c.Delete(); err != nil {
		t.Errorf("c.Delete(): %s\n%s", err, user)
	}
	newcl, err := GetConsoles(user)
	if err != nil {
		t.Errorf("GetConsoles(user:%s): %s", user, err)
	}
	ncl := len(newcl)
	if ncl != lc-1 {
		glog.Errorf("Length of GetConsoles did not go down %v <=> %v", lc, ncl)
	}
}
func TestGetConsoles(t *testing.T) {
	print("GetConsoles()\n")
	seedConsole()
	user := gu(t)
	cl, err := GetConsoles(user)
	if err != nil {
		t.Errorf("GetConsoles(1): %s", err)
	}
	lc := len(cl)
	if lc != 31 {
		t.Errorf("Len GetConsoles(1) 31 <=> %v", lc)
	}
}
func TestConsole_Games(t *testing.T) {
	print("Console.Games()\n")
	c := gsc(t)
	gl, err := c.Games()
	if err != nil {
		t.Errorf("c.Games(): %s", err)
	}
	if len(gl) != 801 {
		t.Errorf("len(c.Games) 801 <=> %v", len(gl))
	}
}
func TestConsole_OwnedGames(t *testing.T) {
	print("Console.OwnedGames()\n")
	print("\tInitial\n")
	seedConsole()
	seedGame()
	user := gu(t)
	c := gsc(t)
	ogc := c.OwnedGames()
	if ogc != 1 {
		t.Errorf("len(c.OwnedGames()) 1 <=> %v", ogc)
	}
	print("\tAfter changing\n")
	g1, err := GetGame(2, user)
	if err != nil {
		t.Errorf("GetGame(2,1): %s", err)
	}
	g1.Has = true
	err = g1.Save()
	if err != nil {
		glog.Errorf("g1.Save(): %s", err)
	}
	ngc := c.OwnedGames()
	if ngc != 2 {
		t.Errorf("len(c.OwnedGames) 2 <=> %v", ngc)
	}
}
func TestConsole_StarContent(t *testing.T) {
	print("Console.StarContent\n")
	c := gsc(t)
	c.Rating = 3
	c.Save()
	u := gu(t)
	c, err := GetConsole("NES", u)
	if err != nil {
		t.Errorf("GetConsole(NES,user: %s", err)
	}
	threestring := template.HTML(" <img id='star_NES_1' src='/static/star_on.png' onclick='set_console_rating(\"NES\",\"NES\",1)'>\n <img id='star_NES_2' src='/static/star_on.png' onclick='set_console_rating(\"NES\",\"NES\",2)'>\n <img id='star_NES_3' src='/static/star_on.png' onclick='set_console_rating(\"NES\",\"NES\",3)'>\n <img id='star_NES_4' src='/static/star_off.png' onclick='set_console_rating(\"NES\",\"NES\",4)'>\n <img id='star_NES_5' src='/static/star_off.png' onclick='set_console_rating(\"NES\",\"NES\",5)'>\n")
	assert.Equal(t, threestring, c.StarContent(), "StarContent")
}
func TestConsole_ShortName(t *testing.T) {
	print("Console.ShortName\n")
	seedConsole()
	u := gu(t)
	c, err := GetConsole("Game Gear", u)
	if err != nil {
		t.Errorf("GetConsole(Game Gear, user): %s", err)
	}
	assert.Equal(t, "GameGear", c.ShortName(), "console.ShortName")
}
func TestConsole_TotalGames(t *testing.T) {
	print("console.TotalGames\n")
	c := gsc(t)
	assert.Equal(t, 801, c.TotalGames(), "TotalGames")
}
func TestConsole_WantedGames(t *testing.T) {
	print("console.WantedGames\n")
	c := gsc(t)
	gl, err := c.Games()
	if err != nil {
		t.Errorf("c.Games(): %s", err)
	}
	print("\tInitial\n")
	wg, err := c.WantedGames()
	if err != nil {
		t.Errorf("c.WantedGames(): %s", err)
	}
	assert.Equal(t, 0, len(wg), "Initial Wanted Games(0)")
	print("\tMarking 10 as wanted\n")
	for i := 0; i <= 10; i++ {
		gl[i].Want = true
		if err := gl[i].Save(); err != nil {
			t.Errorf("gl[%v].Save(): %s", i, err)
		}
	}
	wg, err = c.WantedGames()
	if err != nil {
		t.Errorf("c.WantedGames(): %s", err)
	}
	assert.Equal(t, 10, len(wg), "WantedGames")
	print("\tMarking console as wantgames\n")
	c.WantGames = true
	if err := c.Save(); err != nil {
		t.Errorf("c.Save(): %s", err)
	}
	wg, err = c.WantedGames()
	if err != nil {
		t.Errorf("c.WantedGames(): %s", err)
	}
	assert.Equal(t, 800, len(wg), "WantedGames(console.WantGames=true)")
}
func initConsole() {
	c, err := goconfig.ReadConfigFile("test_config")
	db_name, err := c.GetString("DB", "db")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_host, err := c.GetString("DB", "host")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_user, err := c.GetString("DB", "user")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db_pass, err := c.GetString("DB", "pass")
	if err != nil {
		err.Error()
		fmt.Println(err)
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@"+db_host+"/"+db_name)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("sh", "seed.sh").Output()
	if err != nil {
		glog.Errorf("sh create_tables.sh %s\n%s", err, out)
		return
	}

	auth.DB(db)
	ConsoleDB(db)
	initGame()
}
func seedConsole() {
	out, err := exec.Command("sh", "seed.sh").Output()
	if err != nil {
		glog.Errorf("sh seed.sh %s\n%s", err, out)
		return
	}
}
func gsc(t *testing.T) (c Console) {
	seedConsole()
	user := gu(t)
	c, err := GetConsole("NES", user)
	if err != nil {
		glog.Errorf("GetConsole(NES,1): %s", err)
	}
	return c
}
func gu(t *testing.T) auth.User {
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("GetUser(1): %s", err)
	}
	return user
}
