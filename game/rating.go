package game
// Rating

import (
	"database/sql"
	"github.com/golang/glog"
	u "github.com/ChrisKaufmann/goutils"
)
var (
	stmtSetRating	*sql.Stmt
	stmtGetRating	*sql.Stmt
)
func ratingDB(d *sql.DB) {
	var err error
	ssr := "replace into ratings (thing_id, rating) values (?,?)"
	stmtSetRating, err = u.Sth(db, ssr)
	if err != nil { glog.Errorf("sth(%s): %s", ssr,err) }
	sgr := "select IFNULL(rating,0) from ratings where thing_id=?"
	stmtGetRating, err = u.Sth(db, sgr)
	if err != nil { glog.Errorf("sth(%s): %s", sgr, err) }
}
func (t Thing) Rating() (r int) {
	err := stmtGetRating.QueryRow(t.ID).Scan(&r)
	switch {
		case err == sql.ErrNoRows:
			return 0
		case err != nil:
			glog.Errorf("stmtGetRating.QueryRow(%s): %s", t.ID, err)
			return 0
	}
	return r
}
func (t Thing) SetRating(r int) (err error) {
	_, err = stmtSetRating.Exec(t.ID,r)
	if err != nil { glog.Errorf("stmtSetRating.Exec(%v,%v): %s", t.ID,r, err) }
	return err
}
