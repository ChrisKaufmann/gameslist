package game

import (
	u "github.com/ChrisKaufmann/goutils"
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestCollection(t *testing.T) {
	print("Testing Collection\n")
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
	ec(t, "collection.Consoles",err)
	vl(t,"console length",3,len(cl))

	//Delete a console
	print("\tDeleting Console\n")
	err = coll.DeleteConsole(c3)
	ec(t, "collection.DeleteConsole",err)
	cl, err = coll.Consoles()
	vl(t,"console length after delete",2,len(cl))

	//Add games
	print("\tAdd games\n")
	err = coll.AddGame(g1)
	ec(t, "Add Game 1", err)
	err = coll.AddGame(g2)
	ec(t, "Add Game 2", err)
	err = coll.AddGame(g3)
	ec(t, "Add Game 3", err)

	//Games
	print("\tcollection.Games()\n")
	gl, err := coll.Games()
	ec(t,"collection.Games",err)
	vl(t,"coll.Games length",3,len(gl))

	//console games
	print("\tConsole Games\n")
	cgl, err := coll.ConsoleGames(c1)
	ec(t, "coll.ConsoleGames", err)
	vl(t, "coll.ConsoleGames.length",2,len(cgl))

	//Delete games
	print("\tDelete game from collection\n")
	err = coll.DeleteGame(g3)
	ec(t, "coll.DeleteGame", err)
	gl, err = coll.Games()
	ec(t, "coll.games", err)
	vl(t, "coll.Games.length",2,len(gl))

}
