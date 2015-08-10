package game
//Consoles - crud for consoles

import (
  "database/sql"
  "fmt"
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
	stmtSaveConsole,err = u.Sth(db, "update consoles set name=? where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
	if c.ID < 1 || c.ID != int(c.ID){
		err = errors.New("Bad ID passed to console.Save()")
	}
	_, err = stmtSaveConsole.Exec(c.Name, c.ID)
	return err
}
func (c Console) Games() (gl []Game, err error) {
	gl, err = GetGamesByConsole(c)
	return gl,err
}
func (c Console) FindGame(n string) (gl []Game,err error) {
	tgl,err := c.Games()
	if err != nil {err.Error();fmt.Println(err);return gl,err}
	for _,g := range tgl {
		if strings.Contains(strings.ToLower(g.Name), strings.ToLower(n)){
			gl=append(gl,g)
		}
	}
	return gl,err
}
func (c Console) Delete() (err error) {
	stmtDeleteConsole,err = u.Sth(db, "delete from consoles where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
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
	stmtAddConsole,err = u.Sth(db, "insert into consoles (name) values(?)")
	if err != nil {err.Error();fmt.Println(err);return c,err}
	c.Name=n
	result, err := stmtAddConsole.Exec(c.Name)
	lid, err := result.LastInsertId()
	c.ID=int(lid)
	return c,err
}
func GetAllConsoles() (cl []Console) {
	var stmt,err =u.Sth(db, "select id, name from consoles where 1")
	if err != nil {err.Error();fmt.Println(err);return cl}
	rows, err := stmt.Query()
	if err != nil {
		err.Error();fmt.Println(err)
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
