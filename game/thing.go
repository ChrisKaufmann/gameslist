package game
// things - games, consoles, boxes, manuals

import (
	"fmt"
	"strings"
	u "github.com/ChrisKaufmann/goutils"
	"errors"
	"github.com/golang/glog"
	"database/sql"
	"github.com/ChrisKaufmann/easymemcache"
)
const thingSelectString = " things.id,IFNULL(things.name,''),IFNULL(things.parent_id,0),IFNULL(things.type,'') "

type Thing struct {
	ID			int
	Name		string
	ParentID	int		//if it's a manual, the game/console it's a manual for, for example
	Type		string	//console,game,manual,box,review
}
var (
	usemc			bool
	stmtAddThing	*sql.Stmt
	stmtGetThing	*sql.Stmt
	stmtSaveThing	*sql.Stmt
	stmtDeleteThing *sql.Stmt
	stmtGetGames	*sql.Stmt
	stmtGetBox		*sql.Stmt
	stmtGetManual	*sql.Stmt
	stmtGetAllConsoles *sql.Stmt
	stmtGetAllGames *sql.Stmt
	stmtGetAllThings *sql.Stmt
	db				*sql.DB
	mc				*easymemcache.Client
)
func DB(d *sql.DB){
	db=d
	var err error
	stmtSaveThing,err =u.Sth(db,"update things set name=?, type=?, parent_id=? where id=? limit 1")
	if err != nil {glog.Fatalf("stmtSaveThing-u.Sth(db,'update things set name=?, type=?, parent_id=? where id=? limit 1'):%s",err)}
	stmtDeleteThing,err = u.Sth(db, "delete from things where id=? limit 1")
	if err != nil {glog.Fatalf("stmtDeleteThing-u.St(db,'delete from things where id=? limit 1'): %s", err)}
	stmtGetGames, err = u.Sth(db, "select "+thingSelectString+"	from things where type='game' and parent_id=? order by things.name ASC")
	if err != nil {glog.Fatalf("stmtGetGames-u.Sth(db'select things.id,IFNULL(things.name,''),IFNULL(things.parent_id,''),things.type from things where type='game' and parent_id=? order by things.name ASC'):%s", err)}
	stmtGetBox, err = u.Sth(db, "select "+thingSelectString+" from things where type='box' and parent_id=?")
	if err != nil {glog.Fatalf("stmtGetBox-u.Sth(db,'select "+thingSelectString+" from things where type='box' and parent_id=?): %s", err)}
	stmtGetManual, err = u.Sth(db, "select "+thingSelectString+" from things where type='manual' and parent_id=?")
	if err != nil {glog.Fatalf("stmtGetManual-u.Sth(db,'select "+thingSelectString+" from things where type='manual' and parent_id=?'): %s", err)}
	stmtAddThing,err = u.Sth(db, "insert into things (name,type) values (?,?)")
	if err != nil {glog.Fatalf("stmtAddThing-u.Sth(db,'insert into things (name,type) values (?,?)'): %s", err)}
	stmtGetAllConsoles, err = u.Sth(db, "select "+thingSelectString+" from things where type='console' order by things.name ASC")
	if err != nil {glog.Fatalf("stmtGetAllConsoles-u.Sth(db,'select "+thingSelectString+" from things where type='console' order by things.name ASC", err)}
	stmtGetAllGames, err = u.Sth(db, "select "+thingSelectString+" from things where type='game' order by things.name ASC")
	if err != nil {glog.Fatalf("stmtGetAllGames-u.Sth(db,'select "+thingSelectString+" from things where type='game' order by things.name ASC'): %s", err)}
	stmtGetThing,err = u.Sth(db,"select "+thingSelectString+" from things where id= ?")
	if err != nil {glog.Fatalf("stmtGetThing-u.Sth(db,'select "+thingSelectString+" from things where id= ?'): %s", err) }
	stmtGetAllThings,err = u.Sth(db,"select "+thingSelectString+" from things where 1")
	if err != nil {glog.Fatalf("stmtGetAllThings-u.Sth(db,'select "+thingSelectString+" from things where 1'): %s", err) }
	ratingDB(db)
	reviewDB(db)
	collectionDB(db)
	conditionDB(db)
}
func MemCache(nmc *easymemcache.Client) () {
	mc=nmc
	usemc=true
}

// object functions
func (t Thing)String() string {
	return fmt.Sprintf("ID: %v,Name: %s, Type: %s, ParentID: %v, ", t.ID, t.Name,t.Type,t.ParentID)
}
func (t Thing)Save() (err error){
	if t.ID < 1 || t.ID != int(t.ID){
		err=errors.New("Bad ID in to thing.Save()")
		return err
	}
	if t.Type == "" {
		glog.Error("thing.Save-No type in thing")
		err=errors.New("no type in thing")
		return err
	}
	_,err = stmtSaveThing.Exec(t.Name,t.Type,t.ParentID,t.ID)
	if err != nil {glog.Errorf("thing.Save()-stmtSaveThing.Exec(%s,%s,%s,%s): %s",t.Name,t.Type,t.ParentID,t.ID, err) }
	return err
}
func (t Thing)Delete() (err error) {
    if t.ID < 1 || t.ID != int(t.ID){
        err=errors.New("Bad ID passed to thing.Delete()")
        return err
    }
	_, err = stmtDeleteThing.Exec(t.ID)
	if err != nil {glog.Errorf("thing.Delete()-stmtDeleteThing.Exec(%s): %s",t.ID,err)}
	return err
}
func (t Thing)Games() (tl []Thing, err error) {
	tl,err = getThingsFromSthP(stmtGetGames, t.ID)
	return tl, err
}
func (t Thing)Box() (nt Thing) {
	tl,_ := getThingsFromSthP(stmtGetBox, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name + " Box", "box")
		if err != nil {glog.Errorf("Box()-AddThing(%s Box, 'box'): %s",t.Name,err);return nt }
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {glog.Errorf("Box()-Save():%s",err)}
		return nt
	}
	return tl[0]
}
func (t Thing)Manual() (nt Thing) {
	tl, _ := getThingsFromSthP(stmtGetManual, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name +" Manual", "manual")
		if err != nil {glog.Errorf("thing.Manual()-AddThing(%s Manual,manual):%s",t.Name,err);return nt}
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {glog.Errorf("thing.Manual()-Save(): %s", err)}
		return nt
	}
	return tl[0]
}
func (console Thing)AddGame(n string)(nt Thing, err error) {
	nt,err = AddThing(n, "game")
	if err != nil {glog.Errorf("AddGame: AddThing(%s,'game'):%s",n,err);return nt, err}
	nt.ParentID=console.ID
	err = nt.Save()
	if err != nil {glog.Errorf("AddGame: nt.Save(): %s", err);}
	if usemc {	mc.DeleteLike("games_") }
	return nt, err
}
//non object functions down here
func Search(ss string) (tl []Thing, err error) {
	stmtSearch, err := u.Sth(db, "select "+thingSelectString+" from things where name like '%"+ss+"%' and type in('game','console');")
	if err != nil {glog.Errorf("stmtSearch.sth (select "+thingSelectString+" from things where name like '%?%') %s", err) ;return tl, err}
	tl, err = getThingsFromSth(stmtSearch)
	if err != nil {glog.Errorf("Search(%s): getThingsFromSthP(stmtSearch,%s): %s", ss,ss,err) }
	return tl, err
}
func GetGame(id int) (t Thing, err error) {
	t,err = GetThing(id)
	if err != nil {glog.Errorf("GetGame(%s)-GetThing(%s): %s",id,id,err)}
	return t,err
}
func GetAllThings() (tl []Thing, err error) {
	tl, err = getThingsFromSth(stmtGetAllThings)
	if err != nil {glog.Errorf("GetAllThings()-getThingsFromSth(stmtGetAllThings): %s", err)}
    return tl,err
}
func AddThing(n string, ty string) (t Thing, err error) {
	t.Name = n
	t.Type = ty
    result, err := stmtAddThing.Exec(t.Name,t.Type)
	if err != nil {glog.Errorf("AddThing(%s,%s)-stmtAddThing.Exec(%s,%s): %s",n,ty,t.Name,t.Type,err)}
	lid, err := result.LastInsertId()
	if err != nil {glog.Errorf("AddThing(%s,%s)-result.LastInsertId(): %s",n,ty,err)}
	t.ID=int(lid)
	return t,err
}
func getThingsFromSthPP(stmt *sql.Stmt,ip interface{},pp interface{}) (tl []Thing, err error) {
	param := u.Tostr(ip)
	param2 := u.Tostr(pp)
	rows, err := stmt.Query(param,param2)
	if err != nil {
		glog.Error("getThingsFromSthPP: stmt.query(%s,%s): %s",param,param2,err)
		return tl, err
	}
	for rows.Next() {
		var t Thing
		rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.Type)
		tl = append(tl, t)
	}
	return tl, err
}
func getThingsFromSthP(stmt *sql.Stmt,ip interface{}) (tl []Thing, err error) {
	param := u.Tostr(ip)
	rows, err := stmt.Query(param)
	if err != nil {
		glog.Errorf("getThingsFromSthP(%s): %s", param, err)
		return tl, err
	}
	for rows.Next() {
		var t Thing
		rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.Type)
		tl = append(tl, t)
	}
	return tl, err
}
func getThingsFromSth(stmt *sql.Stmt) (tl []Thing, err error) {
	rows, err := stmt.Query()
	if err != nil {
		glog.Errorf("getThingsFromSth(): %s", err)
		return tl, err
	}
	for rows.Next() {
		var t Thing
		rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.Type)
		tl = append(tl, t)
	}
	return tl, err
}
func getThingsByParam(p string) (tl []Thing, err error) {
	query := "select "+thingSelectString+" from things where "+p
	stmt, err  := u.Sth(db,query)
	if err != nil {glog.Errorf("getThingsByParam-u.Sth(db,%s): %s", query, err) }
	tl, err = getThingsFromSth(stmt)
	if err != nil {glog.Errorf("getThingsByParam-getThingsFromSth(stmt): %s", err) }
	return tl, err
}
func GetThings(idls []int) (tl []Thing, err error) {
	var idl []string
	for _,i := range(idls) {
		idl = append(idl,u.Tostr(i))
	}
	stmtGet, err := u.Sth(db,"select "+thingSelectString+" from things where id in ("+ strings.Join(idl, ",") +")")
	if err != nil {glog.Errorf("GetThings-stmtGet-u.Sth('select "+thingSelectString+" from things where id in ("+ strings.Join(idl, ",") +"): %s", err) }
	tl, err = getThingsFromSth(stmtGet)
	if err != nil {glog.Errorf("getThingsFromSth(stmtGet): %s", err)}
	return tl, err
}
func GetThing(id interface{})(t Thing, err error) {
	id=u.Tostr(id)
	if id == "" {
		return t,err
	}
	tl,err := getThingsFromSthP(stmtGetThing,id)
	if err != nil {
		glog.Errorf("GetThing:getThingsFromSthP(stmtGetThing,%s): %s", id, err)
		return t, err
	}
	if len(tl) > 0 {
		t = tl[0]
	}
	return t,err
}
func GetAllGames() (gl []Thing, err error) {
	ag,err := getThingsFromSth(stmtGetAllGames)
	if err != nil {glog.Errorf("GetAllGames()-getThingsFromSth(stmtGetAllGames): %s", err) }
	return ag, err
}
func GetAllConsoles() (cl []Thing, err error) {
	ac, err := getThingsFromSth(stmtGetAllConsoles)
	if err != nil {glog.Errorf("GetAllConsoles()-getThingsFromSth(stmtGetAllConsoles): %s",err) }
	return ac, err
}
