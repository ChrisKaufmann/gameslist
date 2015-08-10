package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"testing"
)

func TestConsole(t *testing.T) {
	print("Testing Console\n")
	initTest(t)
//adding consoles
	print("\tAdding consoles\n")
	c1,err := AddConsole("console1")
	_ ,err = AddConsole("console2")
	_ ,err = AddConsole("console3")
	_ ,err = AddConsole("console4")
	c5,err := AddConsole("console5")

//Modify a console
	print("\tModifying a console\n")
	c1.Name="Console A"
	err = c1.Save();ec(t,"Save console", err)
//Get a console
	print("\tGet a console\n")
	tid := c1.ID
	c999 := GetConsole(tid)
	if c999.Name != "Console A" {
		t.Errorf("name does not match for c1 and c999 "+c1.Name+" <=> "+c999.Name,err)
	}
//Delete a console
	print("\tDeleting a console\n")
	err = c5.Delete()
	ec(t,"Delete console", err)

//Get the console list, should be 4 left
	print("\tGetting console list\n")
	cl := GetAllConsoles()
	if len(cl) != 4 {
		t.Errorf("Length of GetAllConsoles list is incorrect, expected 4, got "+u.Tostr(len(cl))+"\n")
	}
//create a couple of games
	print("\tCreating games\n")
	_,err=AddGame("Game1",c1)
	_,err=AddGame("Game2",c1)
	_,err=AddGame("Game3",c1)
	_,err=AddGame("Game4",c1)
	_,err=AddGame("Game5",c1)

//get list of games by console
	print("\tGetting list of games by console\n")
	gl,err := c1.Games()
	ec(t,"Get list by console", err)
	if len(gl)!= 5 {
		t.Errorf("Length of console.Games list is incorrect, expected 5, got "+u.Tostr(len(gl))+"\n")
	}
//find consoles
	print("\tFinding consoles\n")
	cfl := FindConsole("onsol")
	if len(cfl) != 4 {
		t.Errorf("Length of FindConsole expected 4, got "+u.Tostr(len(cfl))+"\n")
	}
	cfl2 := FindConsole("onsole3")
	if len(cfl2) != 1 {
		t.Errorf("Length of Findconsole expected 1, got "+u.Tostr(len(cfl2))+"\n")
	}
//find console game(s)
	print("\tFinding console games\n")
	gfl,err := c1.FindGame("game")
	ec(t,"find console games",err)
	if len(gfl) != 5 {
		t.Errorf("Length of c1.FindGame exptected 4, got "+u.Tostr(len(gfl))+"\n")
	}
	gfl2,err := c1.FindGame("game2")
	ec(t,"findgame", err)
	if len(gfl2) != 1 {
		t.Errorf("Length of c1.FindGame exptected 1, got "+u.Tostr(len(gfl2))+"\n")
	}

}
