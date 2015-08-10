package game

import (
	"fmt"
	"errors"
    u "github.com/ChrisKaufmann/goutils"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/msbranco/goconfig"
    "testing"
)

func TestGame(t *testing.T) {
	print("Testing Game\n")
	initTest(t)

//adding consoles
    print("\tAdding consoles\n")
    c1,err := AddConsole("console1")
	ec(t,"add console",err)
    c2 ,err := AddConsole("console2")
	ec(t,"add console",err)
    _ ,err = AddConsole("console3")
	ec(t,"add console",err)
    _ ,err = AddConsole("console4")
	ec(t,"add console",err)
    c5,err := AddConsole("console5")
	ec(t,"add console",err)
	c1.Name=("Console 1");c1.Save()
	ec(t,"save console 1",err)
	c5.Name=("C5");c5.Save()
	ec(t,"save C5", err)
//add games
	print("\tAdding games\n")
	g1 ,err := AddGame("game1",c1);ec(t,"Add game1",err)
	_ ,err = AddGame("game2",c1);ec(t,"Add game2",err)
	_ ,err = AddGame("game3",c2);ec(t,"Add game3",err)
	_ ,err = AddGame("game4",c2);ec(t,"Add game4",err)
	g5 ,err := AddGame("game5",c1);ec(t,"Add game5",err)
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
	if len(gl) != 4 {
		t.Errorf("Length of GetAllGames incorrect, expected 4, got "+u.Tostr(len(gl))+"\n")
	}

//get games by console - already tested in console, but why not
	print("\tGetting a list of games by console\n")
	gl1, err := GetGamesByConsole(c1)
	ec(t, "Get games by console",err)
	if len(gl1) != 2 {
		t.Errorf("Length of getgamesbyconsole incorrect, expected 2, got "+u.Tostr(len(gl1))+"\n")
	}
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
    _,err = db.Query("drop table if exists consoles")
    ec(t,"drop table consoles", err)
    _,err = db.Query("CREATE TABLE consoles (     id int(10) unsigned NOT NULL AUTO_INCREMENT,      name varchar(255) DEFAULT NULL,       PRIMARY KEY (id)        ) AUTO_INCREMENT=1")
    ec(t,"create table consoles",err)

    _,err = db.Query("Drop table if exists games")
    ec(t,"drop table games", err)
    _,err = db.Query("CREATE TABLE games (    id int(10) unsigned NOT NULL AUTO_INCREMENT,      name varchar(255) NOT NULL,       console_id int(10) unsigned NOT NULL,         PRIMARY KEY (id)            ) AUTO_INCREMENT=1;")
    ec(t,"create table games", err)

	_,err = db.Query("Drop table if exists gamecollection")
	ec(t,"drop table gamecollection",err)
	_,err = db.Query("create table gamecollection ( id int unsigned not null primary key auto_increment, user_id int unsigned not null, game_id int unsigned not null) AUTO_INCREMENT=1;")
	ec(t,"create table gamecollection",err)

	_,err = db.Query("Drop table if exists consolecollection")
	ec(t,"drop table consolecollection",err)
	_,err = db.Query("create table consolecollection ( id int unsigned not null primary key auto_increment, user_id int unsigned not null, console_id int unsigned not null) AUTO_INCREMENT=1;")
	ec(t,"create table consolecollection",err)
}
