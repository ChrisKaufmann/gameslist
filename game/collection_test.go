package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"testing"
)

func TestCollection(t *testing.T) {
	print("Testing Collection\n")
	initTest(t)

	//Get new collection
	print("\tGet new collection\n")
	coll,err := GetCollection(1)
	if coll.UserID != 1 {
		t.Errorf("Collection userid doesn't match, expected 1, got " + u.Tostr(coll.UserID) + "\n")
	}
	//Add some things :)
	c1, err:= AddThing("Console 1", "console");if err!=nil {t.Errorf("AddThing(console 1)", err)}
	c2, err:= AddThing("Console 2", "console");if err!=nil {t.Errorf("AddThing(console 2)", err)}
	c3, err:= AddThing("Console 3", "console");if err!=nil {t.Errorf("AddThing(console 3)", err)}
	g1, err := AddThing("Game1", "game");if err!=nil {t.Errorf("AddThing(game 1)", err)}
	g1.ParentID=c1.ID;g1.Save();if err!=nil {t.Errorf("SaveThing(g1)", err)}
	g2, err := AddThing("Game2", "game");if err!=nil {t.Errorf("AddThing(game 2)", err)}
	g2.ParentID=c1.ID;g2.Save();if err!=nil {t.Errorf("SaveThing(g2)", err)}
	g3, err := AddThing("Game3", "game");if err!=nil {t.Errorf("AddThing(game 3)", err)}
	g3.ParentID=c1.ID;g3.Save();if err!=nil {t.Errorf("SaveThing(g3)", err)}
	b1, err := AddThing("BoxGame1","box");if err!=nil {t.Errorf("AddThing(boxgame1)", err)}
	b1.ParentID=g1.ID;err = b1.Save();if err!=nil {t.Errorf("SaveThing(boxgame1)", err)}
	b2, err := AddThing("BoxGame2","box");if err!=nil {t.Errorf("AddThing(boxgame2)", err)}
	b2.ParentID=g2.ID;err = b2.Save();if err!=nil {t.Errorf("SaveThing(boxgame2)", err)}
	m1, err := AddThing("ManualGame1","manual");if err!=nil {t.Errorf("AddThing(manualgame1)", err)}
	m1.ParentID=g1.ID;err = m1.Save();if err!=nil {t.Errorf("SaveThing(ManualGame1)", err)}

	//Add consoles to collection
	err = coll.Add(c1);if err!=nil {t.Errorf("coll.add(c1)", err)}
	err = coll.Add(c2);if err!=nil {t.Errorf("coll.add(c2)", err)}
	err = coll.Add(c3);if err!=nil {t.Errorf("coll.add(c3)", err)}

	//Get collected consoles
	print("\tGet list of consoles\n")
	cl, err := coll.Consoles();if err!=nil {t.Errorf("coll.Consoles()", err)}
	vl(t,"collection.consoles length",3,len(cl))

	//Delete a console
	print("\tDeleting Console\n")
	err = coll.Delete(c3);if err!=nil {t.Errorf("coll.Delete(c3)", err)}
	cl, err = coll.Consoles();if err!=nil {t.Errorf("coll.Consoles", err)}
	vl(t,"console length after delete",2,len(cl))

	//Add games
	print("\tAdd games\n")
	err = coll.Add(g1);if err!=nil {t.Errorf("coll.add(g1)", err)}
	err = coll.Add(g2);if err!=nil {t.Errorf("coll.add(g2)", err)}
	err = coll.Add(g3);if err!=nil {t.Errorf("coll.add(g3)", err)}

	//Games
	print("\tcollection.Games()\n")
	gl, err := coll.Games();if err!=nil {t.Errorf("coll.Games()", err)}
	vl(t,"coll.Games length",3,len(gl))

	//console games
	print("\tConsole Games\n")
	cgl, err := coll.ConsoleGames(c1)
	ec(t, "coll.ConsoleGames", err)
	vl(t, "coll.ConsoleGames.length",3,len(cgl))

	//Delete games
	print("\tDelete game from collection\n")
	err = coll.Delete(g3)
	ec(t, "coll.DeleteGame", err)
	gl, err = coll.Games()
	ec(t, "coll.games", err)
	vl(t, "coll.Games.length",2,len(gl))

	//Add Manual
	print("\tAdd manual\n")
	err = coll.Add(m1);;if err!=nil {t.Errorf("coll.Add(m1)", err)}

	//Get manuals
	print("\tcollection.Manuals()\n")
	ml, err := coll.Manuals();if err!=nil {t.Errorf("coll.Manuals()", err)}
	vl(t,"coll.Manuals length",1, len(ml))

	//Add Box
	print("\tAdd box\n")
	err = coll.Add(b1);if err!=nil {t.Errorf("coll.Ad(b1)", err)}

	//Get Boxes
	print("\tcollection.Boxes()\n")
	bl, err := coll.Boxes();if err!=nil {t.Errorf("coll.Boxes()", err)}
	vl(t,"coll.Boxes length",1,len(bl))

	//Test to see if get thing in collection
	print("\tcollection.Have()\n")
	hg1 := coll.Have(g1);if hg1 == false || err != nil{t.Errorf("coll.Have(g1)", err)}
	print("\tcoll.Have(c3)\n")
	hc3 := coll.Have(c3);if hc3 == true || err != nil{t.Errorf("coll.Have(c3)", err)}

	coll.Print()
}
