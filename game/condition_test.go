package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"testing"
)

func TestConditions(t *testing.T) {
	print("Testing Condition\n")
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

	//Add consoles to collection
	err = coll.Add(c1);if err!=nil {t.Errorf("coll.add(c1): %s", err)}
	err = coll.Add(c2);if err!=nil {t.Errorf("coll.add(c2): %s", err)}

	//GetMyThings for the consoles
	mc1 := coll.GetMyThing(c1)
	mc2 := coll.GetMyThing(c2)

	//Check for 0 for ratings that shouldn't exist yet
	print("\tInitial Conditions\n")
	mc1r := mc1.Condition()
	if mc1r != 0 {t.Errorf("Initial rating for mc1: expected 0 got %v", mc1r)}
	mc2r := mc2.Condition()
	if mc2r != 0 {t.Errorf("Initial rating for mc2: expected 0 got %v", mc2r)}

	//Set ratings and verify
	print("\tSet Conditions\n")
	err = mc1.SetCondition(3)
	if err != nil {t.Errorf("mc1.SetCondition %s", err)}
	mc1r2 := mc1.Condition()
	if mc1r2 != 3 {t.Errorf("mc1 Condition, expected 3 got %v",mc1r2)}
	err = mc2.SetCondition(5)
	if err != nil {t.Errorf("mc2.SetCondition %s", err)}
	mc2r2 := mc2.Condition()
	if mc2r2 != 5 {t.Errorf("mc1 Condition, expected 5 got %v",mc1r2)}

	//GetAllConditions and compare to what we know reality is
	print("\tGetAllConditions\n")
	rl := coll.GetConditions()
	if rl[mc1.ID] != 3 {t.Errorf("GetConditions, mc1 expected 3 got %v", rl[mc1.ID])}
	if rl[mc2.ID] != 5 {t.Errorf("GetConditions, mc2 expected 5 got %v", rl[mc2.ID])}

}
