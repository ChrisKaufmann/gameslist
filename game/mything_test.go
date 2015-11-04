package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"testing"
)

func TestMyThings(t *testing.T) {
	print("Testing MyThing\n")
	initTest(t)

	//Get new collection
	print("\tGet new collection\n")
	coll,err := GetCollection(1)
	if coll.UserID != 1 {
		t.Errorf("Collection userid doesn't match, expected 1, got " + u.Tostr(coll.UserID) + "\n")
	}
	//Add some things :)
	c1, err:= AddThing("Console 1", "console");if err!=nil {t.Errorf("AddThing(console 1): %s", err)}
	c2, err:= AddThing("Console 2", "console");if err!=nil {t.Errorf("AddThing(console 2): %s", err)}
	c3, err:= AddThing("Console 3", "console")

	//Add consoles to collection
	err = coll.Add(c1);if err!=nil {t.Errorf("coll.add(c1): %s", err)}
	err = coll.Add(c2);if err!=nil {t.Errorf("coll.add(c2): %s", err)}

	//GetMyThings for the consoles
	mc1 := coll.GetMyThing(c1)
	//mc2 := coll.GetMyThing(c2)
	mc3 := coll.GetMyThing(c3)

	//Background
	print("\tBackground\n")
	if mc1.Background() != "#c0c0c0" {t.Errorf("mc1.Background, expected '#c0c0c0', got %s", mc1.Background())}
	if mc3.Background() != "white" {t.Errorf("mc3.Background, expected 'white', got %s", mc3.Background())}

	//BoxBackground
	print("\tBoxBackground\n")
	if mc1.BoxBackground() != "white" {t.Errorf("BoxBackground, expected 'white', got %s", mc1.BoxBackground())}
	coll.Add(mc1.Box())
	mc1 = coll.GetMyThing(c1) //refresh mc1
	if mc1.BoxBackground() != "#c0c0c0" {t.Errorf("BoxBackground, expected '#c0c0c0', got %s", mc1.BoxBackground())}

	//ManualBackground
	print("\tManualBackground\n")
	if mc1.ManualBackground() != "white" {t.Errorf("ManualBackground, expected 'white', got %s", mc1.ManualBackground())}
	coll.Add(mc1.Manual())
	mc1 = coll.GetMyThing(c1) //refresh mc1
	if mc1.ManualBackground() != "#c0c0c0" {t.Errorf("ManualBackground, expected '#c0c0c0', got %s", mc1.ManualBackground())}

	//Checked
	print("\tChecked\n")
	if mc1.Checked() != "checked" {t.Errorf("mc1.Checked, expected 'checked', got %s", mc1.Checked())}
	if mc3.Checked() != "unchecked" {t.Errorf("mc3.Checked, expected 'unchecked', got %s", mc3.Checked())}

	//ManualChecked
	print("\tManualChecked\n")
	if mc1.ManualChecked() != "checked" {t.Errorf("mc1.ManualChecked, expected 'checked', got %s",mc1.ManualChecked())}
	if mc3.ManualChecked() != "unchecked" {t.Errorf("mc3.ManualChecked, expected 'unchecked', got %s",mc3.ManualChecked())}

	//boxChecked
	print("\tBoxChecked\n")
	if mc1.BoxChecked() != "checked" {t.Errorf("mc1.BoxChecked, expected 'checked', got %s",mc1.BoxChecked())}
	if mc3.BoxChecked() != "unchecked" {t.Errorf("mc3.BoxChecked, expected 'unchecked', got %s",mc3.BoxChecked())}

	//ConsoleTotal
	print("\tConsoleTotal\n")
	if mc1.ConsoleTotal() != 0 {t.Errorf("ConsoleTotal: expected 0, got %v", mc1.ConsoleTotal())}
	g1,err := c1.AddGame("g1")
	g2,err := c1.AddGame("g2")
	_,err = c1.AddGame("g3")
	if mc1.ConsoleTotal() != 3 {t.Errorf("ConsoleTotal: expected 3, got %v", mc1.ConsoleTotal())}

	//ConsoleOwned
	print("\tConsoleOwned\n")
	if mc1.ConsoleOwned() != 0 {t.Errorf("ConsoleOwned: expected 0, got %v", mc1.ConsoleOwned())}
	coll.Add(g1)
	coll.Add(g2)
	if mc1.ConsoleOwned() != 2 {t.Errorf("ConsoleOwned: expected 2, got %v", mc1.ConsoleOwned())}

}
