package game

import (
	u "github.com/ChrisKaufmann/goutils"
	"fmt"
)

var (
)
type Collection struct {
	UserID		int
}


// Object function
func (coll Collection) Print() {
print("Printing collection\n")
	cons, err := coll.Consoles()
print("len of consoles:"+u.Tostr(len(cons))+"\n")
	if err != nil {err.Error();fmt.Println(err);return}
	for _, c := range cons {
		print(c.Name+"\n")
		cgl, err := coll.ConsoleGames(c)
		if err != nil {err.Error();fmt.Println(err);return}
		for _, g := range cgl {
			print("\t"+g.Name+"\n")
		}
	}
	print("\n")
}
func (coll Collection) Have (t Thing) (h bool) {
	h=false
	stmt, err := u.Sth(db,"select count(*) from collection where thing_id=? and user_id=?")
	if err != nil {err.Error();fmt.Println(err);return h}
	var c int
	err = stmt.QueryRow(t.ID,coll.UserID).Scan(&c)
	if err != nil {err.Error();fmt.Println(err);return h}
	if c > 0 {
		return true
	}
	return false
}
func (coll Collection) Add(t Thing) (err error) {
	stmt,err := u.Sth(db,"insert into collection (user_id, thing_id) values (?,?)")
	if err != nil {err.Error();fmt.Println(err);return err}
	_, err = stmt.Exec(coll.UserID,t.ID)
	return err
}
func (coll Collection) Delete(t Thing) (err error) {
	stmt, err := u.Sth(db, "delete from collection where user_id=? and id=? limit 1")
	_, err = stmt.Exec(coll.UserID, t.ID)
	return err
}
// Console stuff
func (coll Collection) Consoles() (cl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='console'")
	return getThingsFromSthP(stmt,coll.UserID)
}
func (coll Collection) ConsoleGames(c Thing) (gl []Thing, err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game' and things.parent_id=?")
	return getThingsFromSthPP(stmt,coll.UserID,c.ID)
}
// Game stuff
func (coll Collection) Games() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select t.id from things as t, collection as c where c.user_id=? and c.thing_id=t.id and t.type='game'")
	return getThingsFromSthP(stmt,coll.UserID)
}

func (coll Collection) Boxes() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select t.id from things as t, collection as c where c.user_id=? and c.thing_id=t.id and t.type='box'")
	return getThingsFromSthP(stmt,coll.UserID)
}
func (coll Collection) Manuals() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select t.id from things as t, collection as c where c.user_id=? and c.thing_id=t.id and t.type='manual'")
	return getThingsFromSthP(stmt,coll.UserID)
}

//Non Object Functions

func GetCollection(uid int)(coll Collection, err error) {
	coll.UserID=uid
	return coll, err
}
