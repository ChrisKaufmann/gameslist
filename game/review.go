package game

// Revirew

import (
	"database/sql"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
)

var (
	stmtSetReview *sql.Stmt
	stmtGetReview *sql.Stmt
	stmtUserReviews *sql.Stmt
)

func reviewDB(d *sql.DB) {
	var err error
	ssr := "replace into reviews (thing_id,user_id, review) values (?,?,?)"
	stmtSetReview, err = u.Sth(db, ssr)
	if err != nil {
		glog.Errorf("sth(%s): %s", ssr, err)
	}
	sgr := "select IFNULL(review,' ') from reviews where thing_id=? and user_id=?"
	stmtGetReview, err = u.Sth(db, sgr)
	if err != nil {
		glog.Errorf("sth(%s): %s", sgr, err)
	}
	sgur := "select thing_id, review from reviews where user_id=?"
	stmtUserReviews, err = u.Sth(db,sgur)
	if err != nil {
		glog.Errorf("sth(%s): %s", sgur, err)
	}
}
func (t MyThing) Review() (r string) {
	err := stmtGetReview.QueryRow(t.ID, t.Coll.UserID).Scan(&r)
	switch {
	case err == sql.ErrNoRows:
		return ""
	case err != nil:
		glog.Errorf("stmtGetReview.QueryRow(%s): %s", t.ID, err)
		return ""
	}
	return r
}
func (t MyThing) SetReview(r string) (err error) {
	_, err = stmtSetReview.Exec(t.ID, t.Coll.UserID, r)
	if err != nil {
		glog.Errorf("stmtSetReview.Exec(%v,%v,%s): %s", t.ID, t.Coll.UserID, r, err)
	}
	return err
}
func (coll Collection) GetReviews() (map[int]string) {
	var rm = make(map[int]string)
	rows, err := stmtUserReviews.Query(coll.UserID)
	switch {
	case err==sql.ErrNoRows:
		return rm
	case err != nil:
		glog.Errorf("stmtUserReviews.Query(%s): %s", coll.UserID, err)
		return rm
	}
	for rows.Next() {
		var r string
		var tid int
		rows.Scan(&tid, &r)
		rm[tid]=r
	}
	return rm
}
