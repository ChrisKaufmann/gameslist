package game
//Games :)

import (
	"fmt"
	"database/sql"
	u "github.com/ChrisKaufmann/goutils"
	"errors"
)

var (
	stmtGetGameList *sql.Stmt
	stmtSaveGame    *sql.Stmt
	stmtAddGame     *sql.Stmt
	stmtGetGame     *sql.Stmt
	stmtDeleteGame     *sql.Stmt
	stmtGetAllGames *sql.Stmt
	db				*sql.DB
)

func init() {
}
func DB(d *sql.DB) {
	db = d
}

type Game struct {
	ID        int
	Name      string
	ConsoleID int
}

// object functions
func (g Game)Print() {
	print("ID:\t"+u.Tostr(g.ID)+"\nName:\t"+g.Name+"\nconsoleID:\t"+u.Tostr(g.ConsoleID)+"\n")
}
func (g Game)Save() (err error){
	stmtSaveGame,err =u.Sth(db,"update games set name=?, console_id=? where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
	if g.ID < 1 || g.ID != int(g.ID){
		err=errors.New("Bad ID in to game.Save()")
		return err
	}
	_,err = stmtSaveGame.Exec(g.Name,g.ConsoleID,g.ID)
	return err
}
func (g Game)Delete() (err error) {
    stmtDeleteGame,err = u.Sth(db, "delete from games where id=? limit 1")
	if err != nil {err.Error();fmt.Println(err);return err}
    if g.ID < 1 || g.ID != int(g.ID){
        err=errors.New("Bad ID passed to game.Delete()")
        return err
    }
	_, err = stmtDeleteGame.Exec(g.ID)
	return err
}

//non object functions down here
func GetAllGames() (gl []Game, err error) {
	stmtGetAllGames,err = u.Sth(db,"select id,name,console_id from games where 1")
	if err != nil {err.Error();fmt.Println(err);return gl, err}
	rows, err := stmtGetAllGames.Query()
    if err != nil {
        err.Error();fmt.Println(err)
        return gl,err
    }
    for rows.Next() {
        var g Game
        rows.Scan(&g.ID,&g.Name,&g.ConsoleID)
        gl = append(gl,g)
    }
    return gl,err
}
func AddGame(n string, c Console) (g Game, err error) {
	g.Name = n
	g.ConsoleID = c.ID
	stmtAddGame,err = u.Sth(db, "insert into games (name,console_id) values (?,?)")
	if err != nil {err.Error();fmt.Println(err);return g,err}
    result, err := stmtAddGame.Exec(g.Name,g.ConsoleID)
	lid, err := result.LastInsertId()
	g.ID=int(lid)
	return g,err
}
func GetGame(id interface{})(g Game, err error) {
	id=u.Tostr(id)
	if id == "" {
		return g,err
	}
	stmtGetGame,err = u.Sth(db,"select id,name,console_id from games where id= ?")
	if err != nil {err.Error();fmt.Println(err);return g,err}
	err = stmtGetGame.QueryRow(id).Scan(&g.ID, &g.Name, &g.ConsoleID)
	if err != nil {
		err.Error();fmt.Println(err)
	}
	return g,err
}
func GetGamesByConsole(c Console) (gl []Game, err error) {
	stmt,err := u.Sth(db, "select id,name, console_id from games where console_id = ?")
	if err != nil {err.Error();fmt.Println(err);return gl,err}
	rows, err := stmt.Query(u.Tostr(c.ID))
	if err != nil {
		err.Error();fmt.Println(err)
		return gl,err
	}
	for rows.Next() {
		var g Game
		rows.Scan(&g.ID, &g.Name, &g.ConsoleID)
		gl = append(gl, g)
	}
	return gl,err
}
