package game
//Consoles - crud for consoles

import (
  "database/sql"
  u "github.com/ChrisKaufmann/goutils"
  "errors"
  "strings"
)
var (
	stmtGetConsoleList		*sql.Stmt
	stmtSaveConsole			*sql.Stmt
	stmtAddConsole			*sql.Stmt
	stmtDeleteConsole		*sql.Stmt
)
func init() {
}
type Console struct {
	ID		int
	Name	string
}

//object functions
func (c Console) Print() {
	print("ID:\t"+u.Tostr(c.ID)+"\nName:\t"+c.Name+"\n")
}
func (c Console) Save() (err error){
	stmtSaveConsole = u.Sth(db, "update consoles set name=? where id=? limit 1")
	if c.ID < 1 || c.ID != int(c.ID){
		err = errors.New("Bad ID passed to console.Save()")
	}
	_, err = stmtSaveConsole.Exec(c.Name, c.ID)
	return err
}
func (c Console) Games() (gl []Game) {
	return GetGamesByConsole(c)
}
func (c Console) FindGame(n string) (gl []Game) {
	for _,g := range c.Games() {
		if strings.Contains(strings.ToLower(g.Name), strings.ToLower(n)){
			gl=append(gl,g)
		}
	}
	return gl
}
func (c Console) Delete() (err error) {
	stmtDeleteConsole = u.Sth(db, "delete from consoles where id=? limit 1")
	if c.ID < 1 || c.ID != int(c.ID){
		err=errors.New("Bad ID passed to console.Delete()")
		return err
	}
	_, err = stmtDeleteConsole.Exec(c.ID)
	return err
}

//non object functions down here
func FindConsole(n string) (cl []Console) {
	for _,c := range GetAllConsoles() {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(n)) {
			cl=append(cl,c)
		}
	}
	return cl
}
func AddConsole(n string) (c Console, err error) {
	stmtAddConsole = u.Sth(db, "insert into consoles (name) values(?)")
	c.Name=n
	result, err := stmtAddConsole.Exec(c.Name)
	lid, err := result.LastInsertId()
	c.ID=int(lid)
	return c,err
}
func GetAllConsoles() []Console {
	var cl []Console
	var stmt =u.Sth(db, "select id, name from consoles where 1")
	rows, err := stmt.Query()
	if err != nil {
		err.Error()
		return cl
	}
	for rows.Next() {
		var c Console
		rows.Scan(&c.ID,&c.Name)
		cl = append(cl,c)
	}
	return cl
}
func GetConsole(id interface{})  Console {
	id=u.Tostr(id)
	if id == "" {
		var c Console
		return c
	}
	for _,c := range GetAllConsoles() {
		if u.Tostr(c.ID) == id {
			return c
		}
	}
	var c Console
	return c
}
