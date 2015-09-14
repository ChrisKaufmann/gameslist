package game
// things - games, consoles, boxes, manuals

import (
	"fmt"
	"strings"
	u "github.com/ChrisKaufmann/goutils"
	"errors"
	"database/sql"
)
const thingSelectString = " things.id,IFNULL(things.name,''),IFNULL(things.parent_id,''),things.type "

type Thing struct {
	ID			int
	Name		string
	ParentID	int		//if it's a manual, the game/console it's a manual for, for example
	Type		string	//console,game,manual,box
}

var (
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
	stmtHaveThing	*sql.Stmt
	db *sql.DB
)
func DB(d *sql.DB){
	db=d
	var err error
	stmtSaveThing,err =u.Sth(db,"update things set name=?, type=?, parent_id=? where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err)}
	stmtDeleteThing,err = u.Sth(db, "delete from things where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetGames, err = u.Sth(db, "select things.id,IFNULL(things.name,''),IFNULL(things.parent_id,''),things.type from things where type='game' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetBox, err = u.Sth(db, "select "+thingSelectString+" from things where type='box' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetManual, err = u.Sth(db, "select "+thingSelectString+" from things where type='manual' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err)}
	stmtAddThing,err = u.Sth(db, "insert into things (name,type) values (?,?)")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetAllConsoles, err = u.Sth(db, "select "+thingSelectString+" from things where type='console'")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetAllGames, err = u.Sth(db, "select "+thingSelectString+" from things where type='game'")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetThing,err = u.Sth(db,"select "+thingSelectString+" from things where id= ?")
	if err != nil {err.Error();fmt.Println(err)}
	stmtGetAllThings,err = u.Sth(db,"select "+thingSelectString+" from things where 1")
	if err != nil {err.Error();fmt.Println(err)}
	stmtHaveThing, err = u.Sth(db,"select count(*) from collection where thing_id=? and user_id=?")
	if err != nil {err.Error();fmt.Println(err)}
}

// object functions
func (t Thing)Print() {
	print("ID:\t"+u.Tostr(t.ID)+"\nName:\t"+t.Name+"\nType:\t"+u.Tostr(t.Type)+"\nParentID: "+u.Tostr(t.ParentID)+"\n")
}
func (t Thing)Save() (err error){
	if t.ID < 1 || t.ID != int(t.ID){
		err=errors.New("Bad ID in to thing.Save()")
		return err
	}
	if t.Type == "" {
		err=errors.New("no type in thing")
		fmt.Println(err)
		err.Error()
		return err
	}
	_,err = stmtSaveThing.Exec(t.Name,t.Type,t.ParentID,t.ID)
	return err
}
func (t Thing)Delete() (err error) {
    if t.ID < 1 || t.ID != int(t.ID){
        err=errors.New("Bad ID passed to thing.Delete()")
        return err
    }
	_, err = stmtDeleteThing.Exec(t.ID)
	return err
}
func (t Thing)Games() (tl []Thing, err error) {
	return getThingsFromSthP(stmtGetGames, t.ID)
}
func (t Thing)Box() (nt Thing) {
	tl,_ := getThingsFromSthP(stmtGetBox, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name + " Box", "box")
		if err != nil {err.Error();fmt.Println(err);return nt}
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {err.Error();fmt.Println(err);return nt}
		return nt
	}
	return tl[0]
}
func (t Thing)Manual() (nt Thing) {
	tl, _ := getThingsFromSthP(stmtGetManual, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name +" Manual", "manual")
		if err != nil {err.Error();fmt.Println(err);return nt}
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {err.Error();fmt.Println(err);return nt}
		return nt
	}
	return tl[0]
}
//non object functions down here
func GetGame(id int) (t Thing, err error) {
	return GetThing(id)
}
func GetAllThings() (tl []Thing, err error) {
	tl, err = getThingsFromSth(stmtGetAllThings)
    return tl,err
}
func AddThing(n string, ty string) (t Thing, err error) {
	t.Name = n
	t.Type = ty
    result, err := stmtAddThing.Exec(t.Name,t.Type)
	lid, err := result.LastInsertId()
	t.ID=int(lid)
	return t,err
}
func getThingsFromSthPP(stmt *sql.Stmt,ip interface{},pp interface{}) (tl []Thing, err error) {
	param := u.Tostr(ip)
	param2 := u.Tostr(pp)
	rows, err := stmt.Query(param,param2)
	if err != nil {err.Error();fmt.Println(err);return tl,err}
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
	if err != nil {err.Error();fmt.Println(err);return tl,err}
	for rows.Next() {
		var t Thing
		rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.Type)
		tl = append(tl, t)
	}
	return tl, err
}
func getThingsFromSth(stmt *sql.Stmt) (tl []Thing, err error) {
	rows, err := stmt.Query()
	if err != nil {err.Error();fmt.Println(err);return tl,err}
	for rows.Next() {
		var t Thing
		rows.Scan(&t.ID, &t.Name, &t.ParentID, &t.Type)
		tl = append(tl, t)
	}
	return tl, err
}
func getThingsByParam(p string) (tl []Thing, err error) {
	query := "select "+thingSelectString+" from things where "+p
//	print("query="+query+"\n")
	stmt, err  := u.Sth(db,query)
	return getThingsFromSth(stmt)
}
func GetThings(idls []int) (tl []Thing, err error) {
	var idl []string
	for _,i := range(idls) {
		idl = append(idl,u.Tostr(i))
	}
	stmtGet, err := u.Sth(db,"select "+thingSelectString+" from things where id in ("+ strings.Join(idl, ",") +")")
	return getThingsFromSth(stmtGet)
}
func GetThing(id interface{})(t Thing, err error) {
	id=u.Tostr(id)
	if id == "" {
		return t,err
	}
	tl,err := getThingsFromSthP(stmtGetThing,id)
	if err != nil {err.Error();fmt.Println(err);return t,err}
	if len(tl) > 0 {
		t = tl[0]
	}
	return t,err
}
func GetAllGames() (gl []Thing, err error) {
	ag,err := getThingsFromSth(stmtGetAllGames)
	return ag, err
}
func GetAllConsoles() (cl []Thing, err error) {
	ac, err := getThingsFromSth(stmtGetAllConsoles)
	return ac, err
}
