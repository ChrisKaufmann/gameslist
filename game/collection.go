package game

// Collection

import (
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"database/sql"
)
var (
	stmtHaveThing *sql.Stmt
	stmtAddToColl *sql.Stmt
	stmtDelFromColl *sql.Stmt
	stmtCollCons	*sql.Stmt
	stmtConGms		*sql.Stmt
	stmtCollGames	*sql.Stmt
	stmtCollBoxes	*sql.Stmt
	stmtCollManuals *sql.Stmt
	stmtCollThings	*sql.Stmt
)

type Collection struct {
	UserID int
}
func collectionDB(d *sql.DB) {
	var err error

	hts := "select count(*) from collection where thing_id=? and user_id=?"
	stmtHaveThing, err = u.Sth(db,hts)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", hts, err) }

	sas := "insert into collection (user_id, thing_id) values (?,?)"
	stmtAddToColl, err = u.Sth(db, sas)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", sas, err) }

	sds := "delete from collection where user_id=? and thing_id=? limit 1"
	stmtDelFromColl, err = u.Sth(db, sds)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", sds, err) }

	sccs := "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='console' order by things.name ASC"
	stmtCollCons, err = u.Sth(db, sccs)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", sccs, err) }

	scgs := "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game' and things.parent_id=? order by things.name ASC"
	stmtConGms, err = u.Sth(db, scgs)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", scgs, err) }

	scolgs := "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='game'"
	stmtCollGames, err = u.Sth(db, scolgs)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", scolgs, err) }

	scolbs := "select "+thingSelectString+" from things, collection  where collection.user_id=? and collection.thing_id=things.id and things.type='box'"
	stmtCollBoxes, err = u.Sth(db, scolbs)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", scolbs, err) }

	scolmans := "select "+thingSelectString+" from things , collection  where collection.user_id=? and collection.thing_id=things.id and things.type='manual'"
	stmtCollManuals, err = u.Sth(db, scolmans)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", scolmans, err) }

	scollts := "select "+thingSelectString+" from things, collection where collection.user_id=? and collection.thing_id=things.id"
	stmtCollThings, err = u.Sth(db, scollts)
	if err != nil {glog.Fatalf("u.Sth(%s): %s", scollts, err) }
}

// Object function
func (coll Collection) Print() {
	print("Printing collection\n")
	cons, err := coll.Consoles()
	if err != nil {
		glog.Errorf("coll.Consoles(): %s", err)
		return
	}
	for _, mc := range cons {
		c, err := GetThing(mc.ID)
		if err != nil {
			glog.Errorf("GetThing(%v): %s", mc.ID, err)
		}
		print(c.Name + "\n")
		cgl, err := coll.ConsoleGames(c)
		if err != nil {
			glog.Errorf("coll.ConsoleGames(%s): %s", c, err)
			return
		}
		for _, g := range cgl {
			print("\t" + g.Name + "\n")
		}
	}
	print("\n")
}
func (coll Collection) Have(t Thing) (h bool) {
	h = false
	var c int
	err := stmtHaveThing.QueryRow(t.ID, coll.UserID).Scan(&c)
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
	_, err = stmtAddToColl.Exec(coll.UserID, t.ID)
	if err != nil {
		glog.Errorf("stmt.Exec(%v,%v): %s", coll.UserID, t.ID, err)
	}
	return err
}
func (coll Collection) Delete(t Thing) (err error) {
	_, err = stmtDelFromColl.Exec(coll.UserID, t.ID)
	if err != nil {
		glog.Errorf("stmt.Exec(%v,%v): %s", coll.UserID, t.ID, err)
	}
	return err
}

// Console stuff
func (coll Collection) Consoles() (tl []Thing, err error) {
	tl, err = getThingsFromSthP(stmtCollCons, coll.UserID)
	if err != nil {glog.Errorf("getThingsFromSthP(stmtCollCons, %s): %s", coll.UserID, err)}
	return tl, err
}
func (coll Collection) MissingConsoles() (ml []Thing, err error) {
	ac, err := GetAllConsoles()
	if err != nil {
		glog.Errorf("GetAllConsoles(): %s", err)
		return ml, err
	}
	mc, err := coll.Consoles()
	if err != nil {
		glog.Errorf("coll(%s).Consoles(): %s", coll, err)
		return ml, err
	}
	mcc := make(map[int]int)
	for _, c := range mc {
		mcc[c.ID] = 1
	}
	for _, c := range ac {
		if _, ok := mcc[c.ID]; !ok {
			ml = append(ml, c)
		}
	}
	return ml, err
}
func (coll Collection) ConsoleGames(c Thing) ([]Thing, error) {
	gl, err := getThingsFromSthPP(stmtConGms, coll.UserID, c.ID)
	if err != nil {glog.Errorf("getThingsFromSthPP(stmtConGms, %s, %s): %s", coll.UserID, c.ID, err) }
	return gl, err
}
func (coll Collection) OrphanGames() (mtl []Thing, err error) {
	gl, err := coll.Games()
	if err != nil {
		glog.Errorf("coll(%s).Games(): %s", coll, err)
		return mtl, err
	}
	cl, err := coll.Consoles()
	if err != nil {
		glog.Errorf("coll(%s).Consoles(): %s", coll, err)
		return mtl, err
	}
	cm := make(map[string]int)
	for _, c := range cl {
		cm[u.Tostr(c.ID)] = 1
	}
	for _, g := range gl {
		if _, ok := cm[u.Tostr(g.ParentID)]; !ok {
			mtl = append(mtl, g)
		}
	}
	return mtl, err
}


// Game stuff
func (coll Collection) Games() ([]Thing, error) {
	tl, err := getThingsFromSthP(stmtCollGames, coll.UserID)
	if err != nil {glog.Errorf("getThingsFromSthP(stmtCollGames, %s): %s", coll.UserID, err)}
	return tl, err
}

func (coll Collection) Boxes() (gl []Thing, err error) {
	gl, err = getThingsFromSthP(stmtCollBoxes, coll.UserID)
	if err != nil {glog.Errorf("getThingsFromSthp(stmtCollBoxes, %s): %s", coll.UserID, err)}
	return gl, err
}
func (coll Collection) Manuals() (gl []Thing, err error) {
	gl, err = getThingsFromSthP(stmtCollManuals, coll.UserID)
	if err != nil {glog.Errorf("getThingsFromSthP(stmtCollManuals, %s: %s", coll.UserID, err)}
	return gl, err
}

//Misc object functions
func (coll Collection) Things() (tl []Thing, err error) {
	tl, err = getThingsFromSthP(stmtCollThings, coll.UserID)
	if err != nil {
		glog.Errorf("coll.Things getThingsFromSthP(stmt,%s): %s", coll.UserID, err)
	}
	return tl, err
}
func (coll Collection) MyThingsHash() map[int]bool {
	ml, err := coll.Things()
	if err != nil {
		glog.Errorf("coll.MyThingsHash()-Things(): %s", err)
	}
	return ThingHash(ml)
}
func (coll Collection) GetMyThings(tl []Thing)(ptl []MyThing) {
    allthings, err := GetAllThings()
    if err != nil {glog.Errorf("GetMyThings-GetAllThings(): %s", err) }
    ph := coll.MyThingsHash()
    tms := make(map[int]int)
    tbs := make(map[int]int)
    ratingshash := coll.GetRatings()
    //make a hash of all manuals and boxes
    for _,t := range allthings {
        if t.Type=="manual" {tms[t.ParentID]=t.ID}
        if t.Type=="box" {tbs[t.ParentID]=t.ID}
    }
    for _, t := range tl {
        var pt MyThing
		pt.Thing=t
		pt.Coll=coll
        if ph[t.ID] {pt.Has = true }
        if ph[tms[t.ID]] {pt.HasManual = true }
        if ph[tbs[t.ID]] {pt.HasBox = true }
        if ratingshash[t.ID] >0 {pt.rating = ratingshash[t.ID]}
        ptl = append(ptl, pt)
    }
    return ptl
}
func (coll Collection)GetMyThing(t Thing)(pt MyThing) {
    var tl []Thing
    tl = append(tl,t)
    ptl := coll.GetMyThings(tl)
    return ptl[0]
}

// helper function
func (coll Collection) MyThingsFromThings(tl []Thing) (mytl []MyThing) {
	for _, t := range tl {
		var mything MyThing
		mything.Thing=t
		mything.Coll=coll
		mytl = append(mytl, mything)
	}
	return mytl
}

//Non Object Functions
func ThingHash(tl []Thing) map[int]bool {
	h := make(map[int]bool)
	for _, t := range tl {
		h[t.ID] = true
	}
	return h
}
func GetCollection(uid int) (coll Collection, err error) {
	coll.UserID = uid
	return coll, err
}
func (coll Collection) ConsoleSelect() (t template.HTML) {
	cl, err := GetAllConsoles()
	var s string
	if err != nil {
		glog.Errorf("err: %s", err)
	}
	for _, c := range cl {
		s = s + "<option value='" + u.Tostr(c.ID) + "'>" + c.Name + "\n"
	}
	return template.HTML(s)
}
