package game

import (
	u "github.com/ChrisKaufmann/goutils"
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestCollection(t *testing.T) {
	print("Testing Game\n")
	initTest(t)

	//adding consoles
	print("\tAdding consoles\n")
	c1, err := AddConsole("console1")
	ec(t, "add console", err)
	c2, err := AddConsole("console2")
	ec(t, "add console", err)
	c3, err := AddConsole("console3")
	ec(t, "add console", err)
	_, err = AddConsole("console4")
	ec(t, "add console", err)
	c5, err := AddConsole("console5")
	ec(t, "add console", err)
	c1.Name = ("Console 1")
	c1.Save()
	ec(t, "save console 1", err)
	c5.Name = ("C5")
	c5.Save()
	ec(t, "save C5", err)
	//add games
	print("\tAdding games\n")
	g1, err := AddGame("game1", c1)
	ec(t, "Add game1", err)
	g2, err := AddGame("game2", c1)
	ec(t, "Add game2", err)
	g3, err := AddGame("game3", c2)
	ec(t, "Add game3", err)
	_, err = AddGame("game4", c2)
	ec(t, "Add game4", err)
	_, err = AddGame("game5", c1)
	ec(t, "Add game5", err)
	//Get new collection
	print("\tGet new collection\n")
	coll := GetCollection(1)
	if coll.UserID != 1 {
		t.Errorf("Collection userid doesn't match, expected 1, got " + u.Tostr(coll.UserID) + "\n")
	}
	//Add a console
	print("\tAdd consoles\n")
	err = coll.AddConsole(c1)
	ec(t, "Add Console 1", err)
	err = coll.AddConsole(c2)
	ec(t, "Add Console 2", err)
	err = coll.AddConsole(c3)
	ec(t, "Add Console 3", err)

	//Get collected consoles
	print("\tGet list of consoles\n")
	cl, err := coll.Consoles()
	print("cl length:"+u.Tostr(len(cl))+"\n")
	ec(t, "collection.Consoles",err)
	coll.Print()

	//Add games
	print("\tAdd games\n")
	err = coll.AddGame(g1)
	ec(t, "Add Game 1", err)
	err = coll.AddGame(g2)
	ec(t, "Add Game 2", err)
	err = coll.AddGame(g3)
	ec(t, "Add Game 3", err)

	//Delete a console

	//Delete games

}
