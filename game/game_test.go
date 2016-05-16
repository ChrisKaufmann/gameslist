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
	initGame()
}
func TestGetGame(t *testing.T) {
	print("GetGame\n")
	seedGame()
	user, err := auth.GetUser(1)
	if err != nil {
		t.Errorf("auth.GetUser(1): %s", err)
	}
	g, err := GetGame(1, user)
	if err != nil {
		t.Errorf("GetGame(1,1): %s", err)
	}
	if g.User.ID != 1 {
		t.Errorf("g.UserID 1 <=> %v", g.User.ID)
	}
	if g.Has != true {
		t.Errorf("g.Has true <=>%v", g.Has)
		dump()
	}
	if g.HasBox != true {
		t.Errorf("g.HasBox <true <=> %v", g.HasBox)
	}
	if g.HasManual != true {
		t.Errorf("g.HasManual true <=> %v", g.HasManual)
	}
	if g.Rating != 3 {
		t.Errorf("g.Rating 3 <=> %v", g.Rating)
	}
	if g.Review != "is good" {
		t.Errorf("g.Review 'is good' <=> %s", g.Review)
	}

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

	out, err := exec.Command("sh", "seed.sh").Output()
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
