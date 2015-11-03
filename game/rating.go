package game

// Rating

import (
	"database/sql"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
)

var (
	stmtSetRating *sql.Stmt
	stmtGetRating *sql.Stmt
	stmtUserRatings	*sql.Stmt
)

func ratingDB(d *sql.DB) {
	var err error
	ssr := "replace into ratings (thing_id, user_id, rating) values (?,?,?)"
	stmtSetRating, err = u.Sth(db, ssr)
	if err != nil {
		glog.Errorf("sth(%s): %s", ssr, err)
	}
	sgr := "select IFNULL(rating,0) from ratings where thing_id=? AND user_id=?"
	stmtGetRating, err = u.Sth(db, sgr)
	if err != nil {
		glog.Errorf("sth(%s): %s", sgr, err)
	}

	sgur := "select thing_id, rating from ratings where user_id=?"
	stmtUserRatings, err = u.Sth(db, sgur)
	if err != nil {
		glog.Errorf("sth(%s): %s", sgur, err)
	}
}
func (t MyThing) Rating() (r int) {
	err := stmtGetRating.QueryRow(t.ID,t.Coll.UserID).Scan(&r)
	switch {
	case err == sql.ErrNoRows:
		return 0
	case err != nil:
		glog.Errorf("stmtGetRating.QueryRow(%s): %s", t.ID, err)
		return 0
	}
	return r
}
func (t MyThing) SetRating(r int) (err error) {
	_, err = stmtSetRating.Exec(t.ID,t.Coll.UserID, r)
	if err != nil {
		glog.Errorf("stmtSetRating.Exec(%v,%v): %s", t.ID, r, err)
	}
	return err
}
func (coll Collection) GetRatings() (map[int]int) {
	var rm = make(map[int]int)
	rows, err := stmtUserRatings.Query(coll.UserID)
	switch {
		case err == sql.ErrNoRows:
			return rm
		case err != nil:
			glog.Errorf("stmtUserRatings(%s): %s", coll.UserID, err)
			return rm
	}
	for rows.Next() {
		var r int
		var tid int
		rows.Scan(&tid, &r)
		rm[tid]=r
	}
	return rm
}
