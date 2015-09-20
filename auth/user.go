package auth
//  user.auth

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID		int
	Email	string
}

//object functions

//Non object functions
func UserExists(email string)(exists bool) {
    var uid int
    err := stmtGetUserID.QueryRow(email).Scan(&uid)
    switch {
        case err == sql.ErrNoRows:
            exists = false
        case err != nil:
            fmt.Println(err);err.Error();exists=false
        default:
            exists = true
    }
    return exists
}
func AddUser(e string)(us User, err error) {
    if UserExists(e) {
        err = stmtGetUserID.QueryRow(e).Scan(&us.ID)
        return us,err
    }
    result, err := stmtInsertUser.Exec(e)
    if err != nil { print("AddUser("+e+"):stmtInsertUser");fmt.Println(err);err.Error();return us,err}
    lid, err := result.LastInsertId()
    us.ID=int(lid)
	us.Email=e
    return us, err
}
func GetUserByEmail(e string)(us User, err error) {
    if !UserExists(e) {
        err=errors.New("User Doesn't exist")
        return us, err
    }
    err = stmtGetUserID.QueryRow(e).Scan(&us.ID)
    if err != nil {print("GetUserByEmail("+e+")stmtGetUserID");fmt.Println(err);err.Error()}
	us.Email=e
    return us, err
}
func GetUserBySession(s string)(us User, err error) {
//	stmtGetUserByHash, err = u.Sth(db, "select user.id, user.email from user, sessions where user.id=sessions.user_id and sessions.session_hash=?")
//	stmtGetUser,err = us.Sth(db, "select user_id from sessions as s where s.session_hash = ?")
	err = stmtGetUserBySession.QueryRow(s).Scan(&us.ID, &us.Email)
	switch {
		case err == sql.ErrNoRows:
			err = errors.New("No valid session")
			return us, err
		case err != nil:
			print("GEtUserBySession("+s+")");fmt.Println(err);err.Error();return us,err
	}
	return us,err
}
func SessionExists(s string)(e bool) {
//	stmtSessionExists, err := u.Sth(db, "select user_id from sessions where session_hash=?"

	var uid int
	err := stmtSessionExists.QueryRow(s).Scan(&uid)
	switch {
		case err == sql.ErrNoRows:
			return false
		case err != nil:
			print("SessionExists("+s+");stmtSessionExists;");fmt.Println(err);err.Error();return false
		default:
			return true
	}
	return e
}
