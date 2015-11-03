package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"testing"
)

func TestReviews(t *testing.T) {
	print("Testing Review\n")
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

	//Check for 0 for reviews that shouldn't exist yet
	print("\tInitial Reviews\n")
	mc1r := mc1.Review()
	if mc1r != "" {t.Errorf("Initial review for mc1: expected empty got %s", mc1r)}
	mc2r := mc2.Review()
	if mc2r != "" {t.Errorf("Initial review for mc2: expected empty got %s", mc2r)}

	//Set reviews and verify
	print("\tSet Reviews\n")
	err = mc1.SetReview("O hai")
	if err != nil {t.Errorf("mc1.SetReview %s", err)}
	mc1r2 := mc1.Review()
	if mc1r2 != "O hai" {t.Errorf("mc1 Review, expected 'O hai' got %s",mc1r2)}
	err = mc2.SetReview("ñ special")
	if err != nil {t.Errorf("mc2.SetReview %s", err)}
	mc2r2 := mc2.Review()
	if mc2r2 != "ñ special" {t.Errorf("mc1 Review, expected 'ñ special' got %v",mc1r2)}

	//GetAllReviews and compare to what we know reality is
	print("\tGetAllReviews\n")
	rl := coll.GetReviews()
	if rl[mc1.ID] != "O hai" {t.Errorf("GetReviews, mc1 expected 'O hai' got %v", rl[mc1.ID])}
	if rl[mc2.ID] != "ñ special" {t.Errorf("GetReviews, mc2 expected 'ñ special' got %v", rl[mc2.ID])}

}
