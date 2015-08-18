package game
// thing_test.go

import (
	"fmt"
	"errors"
    u "github.com/ChrisKaufmann/goutils"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/msbranco/goconfig"
    "testing"
)

func TestThing(t *testing.T) {
return
	print("Testing Thing\n")
	initTest(t)
//Add things
	print("\tAdd Things\n")
    c1, err:= AddThing("Console 1", "console")
    ec(t, "Add Console 1", err)
    c2, err:= AddThing("Console 2", "console")
    ec(t, "Add Console 2", err)
    _, err= AddThing("Console 3", "console")
    ec(t, "Add Console 3", err)

    g1, err := AddThing("Game1", "game")
    ec(t, "Add Game 1", err)
	g1.ParentID=c1.ID;err = g1.Save()

    g2, err := AddThing("Game2", "game")
    ec(t, "Add Game 2", err)
	g2.ParentID=c1.ID;err = g2.Save()

    g3, err := AddThing("Game3", "game")
    ec(t, "Add Game 3", err)
	g3.ParentID=c1.ID;err = g3.Save()

    g4, err := AddThing("Game4", "game")
    ec(t, "Add Game 4", err)
	g4.ParentID=c2.ID;err = g4.Save()

    g5, err := AddThing("Game5", "game")
    ec(t, "Add Game 5", err)
	g5.ParentID=c2.ID;err = g5.Save()

//save game
	print("\tSaving a game\n")
	g1.Name="Game 1"
	err = g1.Save();ec(t,"Save Game1",err)
//Get a game
	print("\tGetting a game\n")
	g999,err := GetGame(g1.ID)
	ec(t,"get game", err)
	if g999.Name != "Game 1" {
		t.Errorf("name doesn't match for g1,g999: "+g1.Name+"<=>"+g999.Name+"\n", err)
	}

//delete game
	print("\tDeleting a game\n")
	err = g5.Delete();ec(t,"Deleting game", err)

//get list of all games
	print("\tGetting alist of all games\n")
	gl,err := GetAllGames()
	ec(t, "get all games", err)
	vl(t, "Length of GetAllGAmes",4,len(gl))

//get games by console
	print("\tGetting a list of games by console\n")
	gl1, err := c1.Games()
	ec(t, "Get games by console",err)
	vl(t, "Length of GetGamesByConsole",3,len(gl1))
//GetThings
	print("\tGet Things (list of ids)\n")
	idl := []int{1,2,3}
	tl1, err := GetThings(idl)
	ec(t, "GetThings", err)
	vl(t, "Length of GetThings", 3, len(tl1))
//add a manual
	print("\tAdd Manual\n")
	m1,err := AddThing("Game 4 Manual", "manual");if err!=nil {t.Errorf("AddThing(game4manual)", err)}
	m1.ParentID=g4.ID
	err = m1.Save();if err!=nil {t.Errorf("m1.Save()", err)}
//Get a manual that was already added, from a game
	print("\tGetManual(game.Manual()\n")
	m2, err := g4.Manual();if err!=nil {t.Errorf("g4.Manual()", err)}
	if m2.Type != "manual" {t.Errorf("m2.Manual, expected 'manual, got: "+m2.Type, err)}
//Get a manual that shouldn't exist
	print("\tGetManual(game.manual()- not pre-existing)\n")
	m3,err := g1.Manual();if err!=nil {t.Errorf("g1.Manual", err)}
	if m3.ParentID != g1.ID {t.Errorf("m3.ParentID, Expected:"+u.Tostr(g1.ID)+", Got:"+u.Tostr(m3.ParentID)+"\n", err)}
//add a box 
	print("\tAdd Box\n")
	b1,err := AddThing("Game 4 Box", "box");if err!=nil {t.Errorf("AddThing(game4box)", err)}
	b1.ParentID=g4.ID
	err = b1.Save();if err!=nil {t.Errorf("b1.Save()", err)}
//Get a manual that was already added, from a game
	print("\tGetBox(game.Box()\n")
	b2, err := g4.Box();if err!=nil {t.Errorf("g4.Box()", err)}
	if b2.Type != "box" {t.Errorf("m2.Box, expected 'box, got: "+b2.Type, err)}
//Get a manual that shouldn't exist
	print("\tGetBox(game.box()- not pre-existing)\n")
	b3,err := g1.Box();if err!=nil {t.Errorf("g1.Box", err)}
	if b3.ParentID != g1.ID {t.Errorf("b3.ParentID, Expected:"+u.Tostr(g1.ID)+", Got:"+u.Tostr(b3.ParentID)+"\n", err)}

}
func vl(t *testing.T,s string, e interface{}, a interface{}) {
	if e != a {
		err := errors.New("expected: "+u.Tostr(e)+" got: "+u.Tostr(a)+"\n")
		t.Errorf(s, err)
	}
}
func ec(t *testing.T,s string, err error) {
    if err != nil {
        t.Errorf(s, err)
    }
}
func initTest(t *testing.T)  {
    c, err := goconfig.ReadConfigFile("config")
    db_name, err := c.GetString("DB", "db")
    if err != nil {
        err.Error();fmt.Println(err)
    }
    db_host, err := c.GetString("DB", "host")
    if err != nil {
        err.Error();fmt.Println(err)
    }
    db_user, err := c.GetString("DB", "user")
    if err != nil {
        err.Error();fmt.Println(err)
    }
    db_pass, err := c.GetString("DB", "pass")
    if err != nil {
        err.Error();fmt.Println(err)
    }
    db, err = sql.Open("mysql", db_user+":"+db_pass+"@"+db_host+"/"+db_name)
    if err != nil {
        panic(err)
    }
    DB(db)
	_,err = db.Query("Drop table if exists things")
	ec(t,"drop table gamecollection",err)
	_,err = db.Query("CREATE TABLE `things` (  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,  `name` varchar(255) NOT NULL,  `parent_id` int(10) unsigned DEFAULT NULL,  `type` enum('console','game','manual','box') DEFAULT NULL,  PRIMARY KEY (`id`)) AUTO_INCREMENT=1;")
	ec(t,"create table things",err)

	_,err = db.Query("Drop table if exists collection")
	ec(t,"drop table collection",err)
	_,err = db.Query("create table collection ( id int unsigned not null primary key auto_increment, user_id int unsigned not null, thing_id int unsigned not null) AUTO_INCREMENT=1;")
	ec(t,"create table collection",err)
}
