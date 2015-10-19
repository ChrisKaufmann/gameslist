package game
// Collection

import (
	u "github.com/ChrisKaufmann/goutils"
	"fmt"
	"html/template"
	"github.com/golang/glog"
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
	if err != nil {
		glog.Errorf("stmtHaveThing.QueryRow(%v,%v): %s", t.ID, coll.UserID, err)
		return h
	}
	if c > 0 {
		return true
	}
	return false
}
func (coll Collection) Add(t Thing) (err error) {
	stmt,err := u.Sth(db,"insert into collection (user_id, thing_id) values (?,?)")
	if err != nil {
		glog.Errorf("u.Sth(): %s", err)
		return err
	}
	_, err = stmt.Exec(coll.UserID,t.ID)
	if err != nil { glog.Errorf("stmt.Exec(%v,%v): %s", coll.UserID, t.ID, err) }
	return err
}
func (coll Collection) Delete(t Thing) (err error) {
	stmt, err := u.Sth(db, "delete from collection where user_id=? and thing_id=? limit 1")
	_, err = stmt.Exec(coll.UserID, t.ID)
	if err != nil {glog.Errorf("stmt.Exec(%v,%v): %s", coll.UserID, t.ID, err) }
	return err
}
// Console stuff
func (coll Collection) Consoles() (tl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='console' order by things.name ASC")
	tl, err = getThingsFromSthP(stmt,coll.UserID)
	return tl, err
}
func (coll Collection) MissingConsoles() (ml []Thing, err error) {
	ac, err := GetAllConsoles()
	if err != nil {return ml, err}
	mc, err := coll.Consoles()
	if err != nil {return ml, err}
	mcc := make(map[int]int)
	for _,c := range mc {
		mcc[c.ID]=1
	}
	for _,c := range ac {
		if _,ok := mcc[c.ID]; !ok {
			ml = append(ml, c)
		}
	}
	return ml, err
}
func (coll Collection) ConsoleGames(c Thing) ([]Thing, error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game' and things.parent_id=? order by things.name ASC")
	gl, err := getThingsFromSthPP(stmt,coll.UserID,c.ID)
	return gl, err
}
func (coll Collection) OrphanGames() (mtl []Thing,err error) {
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
func (coll Collection) Games() ( []Thing, error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game'")
	tl, err := getThingsFromSthP(stmt,coll.UserID)
	return tl,err
}

func (coll Collection) Boxes() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things, collection  where collection.user_id=? and collection.thing_id=things.id and things.type='box'")
	return getThingsFromSthP(stmt,coll.UserID)
}
func (coll Collection) Manuals() (gl []Thing, err error) {
	stmt, err := u.Sth(db,"select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='manual'")
	return getThingsFromSthP(stmt,coll.UserID)
}

//Misc object functions
func (coll Collection) Things() (tl []Thing,err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things, collection where collection.user_id=? and collection.thing_id=things.id")
	if err != nil { glog.Errorf("coll.Things sth: %s", err); return}
	tl,err= getThingsFromSthP(stmt,coll.UserID)
	if err!= nil {glog.Errorf("coll.Things getThingsFromSthP(stmt,%s): %s", coll.UserID, err)}
	return tl, err
}
func (coll Collection) MyThingsHash() (map[int]bool) {
	ml, err := coll.Things()
	if err != nil {glog.Errorf("coll.MyThingsHash()-Things(): %s", err) }
	return ThingHash(ml)
}

// helper function
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
func ThingHash(tl []Thing)(map[int]bool) {
	h := make(map[int]bool)
	for _,t := range tl {
		h[t.ID]=true
	}
	return h
}
func GetCollection(uid int)(coll Collection, err error) {
	coll.UserID=uid
	return coll, err
}
func (coll Collection) ConsoleSelect() (t template.HTML) {
    cl, err := GetAllConsoles()
	var s string
    if err != nil {
		glog.Errorf("err: %s", err)
    }
    for _,c := range cl {
        s = s+"<option value='"+u.Tostr(c.ID)+"'>"+c.Name+"\n"
    }
    return template.HTML(s)
}
