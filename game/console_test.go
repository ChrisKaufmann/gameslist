package game

import (
	"database/sql"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/msbranco/goconfig"
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
	if c.User.ID != 1 {
		t.Errorf("c.UserID 1 <=> %v", c.User.ID)
	}
	if c.Has != true {
		t.Errorf("c.Has true <=>%v", c.Has)
	}
	if c.HasBox != true {
		t.Errorf("c.HasBox <true <=> %v", c.HasBox)
	}
	if c.HasManual != true {
		t.Errorf("c.HasManual true <=> %v", c.HasManual)
	}
	if c.Rating != 3 {
		t.Errorf("c.Rating 3 <=> %v", c.Rating)
	}
	if c.Review != "is good" {
		t.Errorf("c.Review 'is good' <=> %s", c.Review)
	}

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
	if d.Has != false {
		t.Errorf("d.Has false <=> %v", d.Has)
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
	if d.Manufacturer != "newman1" {
		t.Errorf("d.Manufacturer newman1<=>%s", d.Manufacturer)
	}
	if d.Year != 2000 {
		t.Errorf("d.Year 2000 <=> %v", d.Year)
	}
}
func TestGetConsoles(t *testing.T) {
	print("TestGetConsoles()\n")
	seedConsole()
	user := gu(t)
	cl, err := GetConsoles(user)
	if err != nil {
		t.Errorf("GetConsoles(1): %s", err)
	}
	lc := len(cl)
	if lc != 13 {
		t.Errorf("Len GetConsoles(1) 13 <=> %v", lc)
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
