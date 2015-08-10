package game

//Games :)

import (
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
		err.Error()
		return
	}
	for _, c := range cc {
		print("\tConsole:" + c.Name)
	}
}
func (cl Collection) AddConsole(c Console) (err error) {
	stmt := u.Sth(db,"insert into consolecollection (user_id,console_id) values (?,?)")
	_, err = stmt.Exec(cl.UserID, c.ID)
	return err
}
func (cl Collection) DeleteConsole(c Console) (err error) {
	err = errors.New("DeleteConsole not implemented yet")
	return err
}
func (cl Collection) Games() (gl []Game, err error) {
	err = errors.New("Games not implemented yet")
	return gl, err
}
func (cl Collection) Consoles() (ncl []Console, err error) {
	stmt := u.Sth(db, "select console_id from consolecollection where user_id=?")
	rows, err := stmt.Query(u.Tostr(cl.UserID))
	if err != nil {
		err.Error()
		return ncl , err
	}
	for rows.Next() {
		var id int
		c := GetConsole(id)
		ncl = append(ncl, c)
	}
	return ncl, err
}
func (cl Collection) ConsoleGames() (gl []Game, err error) {
	err = errors.New("ConsoleGames not implemented yet")
	return gl, err
}
func (cl Collection) AddGame(g Game) (err error) {
	stmt := u.Sth(db, "insert into gamecollection (user_id,game_id) values (?,?)")
	_, err = stmt.Exec(cl.UserID, g.ID)
	return err
}
func (cl Collection) DeleteGame(g Game) (err error) {
	stmt := u.Sth(db, "delete from gamecollection where user_id=? and game_id=? limit 1")
	_, err = stmt.Exec(cl.UserID, g.ID)
	return err
}

//non object functions down here
func GetCollection(uid int) (c Collection) {
	c.UserID = int(uid)
	return c
}
