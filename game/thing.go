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
var db *sql.DB

type Thing struct {
	ID			int
	Name		string
	ParentID	int		//if it's a manual, the game/console it's a manual for, for example
	Type		string	//console,game,manual,box
}

func DB(d *sql.DB){
	db=d
}

// object functions
func (t Thing)Print() {
	print("ID:\t"+u.Tostr(t.ID)+"\nName:\t"+t.Name+"\nType:\t"+u.Tostr(t.Type)+"\nParentID: "+u.Tostr(t.ParentID)+"\n")
}
func (t Thing)Save() (err error){
	stmtSave,err :=u.Sth(db,"update things set name=?, type=?, parent_id=? where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
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
	_,err = stmtSave.Exec(t.Name,t.Type,t.ParentID,t.ID)
	return err
}
func (t Thing)Delete() (err error) {
    stmtDelete,err := u.Sth(db, "delete from things where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
    if t.ID < 1 || t.ID != int(t.ID){
        err=errors.New("Bad ID passed to thing.Delete()")
        return err
    }
	_, err = stmtDelete.Exec(t.ID)
	return err
}
func (t Thing)Games() (tl []Thing, err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things where type='game' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err);return tl,err}
	return getThingsFromSthP(stmt, t.ID)
}
func (t Thing)Box() (nt Thing,err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things where type='box' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err);return nt,err}
	tl, err := getThingsFromSthP(stmt, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name + " Box", "box")
		if err != nil {err.Error();fmt.Println(err);return nt,err}
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {err.Error();fmt.Println(err);return nt,err}
		return nt, err
	}
	return tl[0], err
}
func (t Thing)Manual() (nt Thing,err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things where type='manual' and parent_id=?")
	if err != nil {err.Error();fmt.Println(err);return nt,err}
	tl, err := getThingsFromSthP(stmt, t.ID)
	if len(tl) < 1 {
		nt, err := AddThing(t.Name +" Manual", "manual")
		if err != nil {err.Error();fmt.Println(err);return nt,err}
		nt.ParentID=t.ID
		err =  nt.Save()
		if err != nil {err.Error();fmt.Println(err);return nt,err}
		return nt, err
	}
	return tl[0], err
}

//non object functions down here
func GetGame(id int) (t Thing, err error) {
	return GetThing(id)
}
func GetAllThings() (tl []Thing, err error) {
	stmtGetAll,err := u.Sth(db,"select "+thingSelectString+" from things where 1")
	tl, err = getThingsFromSth(stmtGetAll)
    return tl,err
}
func AddThing(n string, ty string) (t Thing, err error) {
	t.Name = n
	t.Type = ty
	stmtAdd,err := u.Sth(db, "insert into things (name,type) values (?,?)")
	if err != nil {err.Error();fmt.Println(err);return t,err}
    result, err := stmtAdd.Exec(t.Name,t.Type)
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
	stmtGet,err := u.Sth(db,"select "+thingSelectString+" from things where id= "+u.Tostr(id))
	if err != nil {err.Error();fmt.Println(err);return t,err}
	tl,err := getThingsFromSth(stmtGet)
	if err != nil {err.Error();fmt.Println(err);return t,err}
	if len(tl) > 0 {
		t = tl[0]
	}
	return t,err
}
func GetAllGames() (gl []Thing, err error) {
	stmt, err := u.Sth(db, "select id, name, type from things where type='game'")
	return getThingsFromSth(stmt)
}
func GetAllConsoles() (cl []Thing, err error) {
	stmt, err := u.Sth(db, "select "+thingSelectString+" from things where type='console'")
	return getThingsFromSth(stmt)
}
