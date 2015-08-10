package game

//Collection :)

import (
	"fmt"
	"errors"
	u "github.com/ChrisKaufmann/goutils"
)

type Collection struct {
	UserID int
}

// object functions
func (cl Collection) Print() {
	cc, err := cl.Consoles()
	if err != nil {
		err.Error();fmt.Println(err)
		return
	}
	print("Print collection\n")
	for _, c := range cc {
		print("\tConsole:" + c.Name +"\n")
		gl, err := cl.ConsoleGames(c)
		if err != nil {
			err.Error();fmt.Println(err)
			return
		}
		for _, g := range gl {
			print("\t\tGame:"+g.Name +"\n")
		}
	}
}
func (cl Collection) AddConsole(c Console) (err error) {
	stmt,err := u.Sth(db,"insert into consolecollection (user_id,console_id) values (?,?)")
	if err != nil {err.Error();fmt.Println(err);return err}
	_, err = stmt.Exec(cl.UserID, c.ID)
	return err
}
func (cl Collection) DeleteConsole(c Console) (err error) {
	stmt, err := u.Sth(db, "delete from consolecollection where user_id=? and console_id=? limit 1")
	if c.ID < 1 || c.ID != int(c.ID){
		err = errors.New("Bad id passed to collection.DeleteConsole")
		return err
	}
	_, err = stmt.Exec(cl.UserID, c.ID)
	return err
}
func (cl Collection) Games() (gl []Game, err error) {
	stmt,err := u.Sth(db, "select game_id from gamecollection where user_id=?")
	if err != nil {err.Error();fmt.Println(err); return gl, err}
	rows, err := stmt.Query(u.Tostr(cl.UserID))
	if err != nil {err.Error();fmt.Println(err); return gl, err}
	for rows.Next() {
		var gid int
		rows.Scan(&gid)
		g,err := GetGame(gid)
		if err != nil {err.Error();fmt.Println(err); return gl, err}
		gl = append(gl,g)
	}
	return gl, err
}
func (cl Collection) Consoles() (ncl []Console, err error) {
	stmt,err := u.Sth(db, "select console_id from consolecollection where user_id=?")
	if err != nil {err.Error();fmt.Println(err);return ncl,err}
	rows, err := stmt.Query(u.Tostr(cl.UserID))
	if err != nil {
		err.Error();fmt.Println(err)
		return ncl , err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		c := GetConsole(id)
		ncl = append(ncl, c)
	}
	return ncl, err
}
func (cl Collection) ConsoleGames(c Console) (gl []Game, err error) {
	stmt,err := u.Sth(db, "select g.id from games as g, consolecollection as cc,gamecollection as gc  where cc.user_id=? and cc.console_id=? and cc.console_id=g.console_id and g.ID=gc.game_id")
	if err != nil {err.Error();fmt.Println(err);return gl,err}
	rows, err := stmt.Query(u.Tostr(cl.UserID),c.ID)
	if err != nil {
		err.Error();fmt.Println(err)
		return gl, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		g,err := GetGame(id)
		if err != nil {err.Error();fmt.Println(err);return gl,err}
		gl = append(gl,g)
	}
	return gl, err
}
func (cl Collection) AddGame(g Game) (err error) {
	stmt,err := u.Sth(db, "insert into gamecollection (user_id,game_id) values (?,?)")
	if err != nil {err.Error();fmt.Println(err);return err}
	_, err = stmt.Exec(cl.UserID, g.ID)
	return err
}
func (cl Collection) DeleteGame(g Game) (err error) {
	stmt,err := u.Sth(db, "delete from gamecollection where user_id=? and game_id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
	_, err = stmt.Exec(cl.UserID, g.ID)
	return err
}

//non object functions down here
func GetCollection(uid int) (c Collection) {
	c.UserID = int(uid)
	return c
}
