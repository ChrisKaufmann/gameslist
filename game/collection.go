package game
// Collection

import (
	u "github.com/ChrisKaufmann/goutils"
	"fmt"
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
	for _, mc := range cons {
		c, err := GetThing(mc.ID)
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
	var c int
	err := stmtHaveThing.QueryRow(t.ID,coll.UserID).Scan(&c)
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
	stmt, err := u.Sth(db, "delete from collection where user_id=? and thing_id=? limit 1")
	_, err = stmt.Exec(coll.UserID, t.ID)
	if err != nil {err.Error();fmt.Println(err);return err}
	print("Deleted "+u.Tostr(t.ID))
	return err
}
// Console stuff
func (coll Collection) Consoles() (mytl []MyThing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='console'")
	tl, err := getThingsFromSthP(stmt,coll.UserID)
	return coll.MyThingsFromThings(tl), err
}
func (coll Collection) MissingConsoles() (mytl []MyThing, err error) {
	ac, err := GetAllConsoles()
	if err != nil {return mytl, err}
	mc, err := coll.Consoles()
	if err != nil {return mytl, err}
	mcc := make(map[int]int)
	for _,c := range mc {
		mcc[c.ID]=1
	}
	var ml []Thing
	for _,c := range ac {
		if _,ok := mcc[c.ID]; !ok {
			ml = append(ml, c)
		}
	}
	mytl=coll.MyThingsFromThings(ml)
	return mytl, err
}
func (coll Collection) ConsoleGames(c Thing) ([]MyThing, error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game' and things.parent_id=?")
	gl, err := getThingsFromSthPP(stmt,coll.UserID,c.ID)
	return coll.MyThingsFromThings(gl), err
}
func (coll Collection) OrphanGames() (mtl []MyThing,err error) {
	gl, err := coll.Games()
	if err != nil {return mtl, err}
	cl, err := coll.Consoles()
	if err != nil {return mtl, err}
	cm := make(map[string]int)
	for _, c := range cl {
		cm[u.Tostr(c.ID)]=1
	}
	for _, g := range gl {
		if _, ok := cm[u.Tostr(g.ParentID)]; !ok {
			mtl = append(mtl, g)
		}
	}
	return mtl, err
}
// Game stuff
func (coll Collection) Games() ( []MyThing, error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game'")
	tl, err := getThingsFromSthP(stmt,coll.UserID)
	return coll.MyThingsFromThings(tl), err
}

func (coll Collection) Boxes() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things, collection  where collection.user_id=? and collection.thing_id=things.id and things.type='box'")
	return getThingsFromSthP(stmt,coll.UserID)
}
func (coll Collection) Manuals() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='manual'")
	return getThingsFromSthP(stmt,coll.UserID)
}
func (coll Collection) MyThingsFromThings(tl []Thing)(mytl []MyThing) {
	ml, err := coll.Manuals()
	if err != nil {err.Error();fmt.Println(err)}
	manual_map := make(map[int]int)
	for _,m := range ml {
		manual_map[m.ParentID]=1
	}
	bm, err := coll.Boxes()
	if err != nil {err.Error();fmt.Println(err)}
	box_map := make(map[int]int)
	for _,b := range bm {
		box_map[b.ParentID]=1
	}
	for _,t := range tl {
		mything := MyThing{t,coll}
		mytl = append(mytl, mything)
	}
	return mytl
}

//Non Object Functions

func GetCollection(uid int)(coll Collection, err error) {
	coll.UserID=uid
	return coll, err
}
