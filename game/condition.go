package game

//Condition

import (
	"database/sql"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
)
var (
	stmtSetCondition	*sql.Stmt
	stmtGetCondition	*sql.Stmt
	stmtUserConditions	*sql.Stmt
)

func conditionDB(d *sql.DB) {
	var err error
	ssc := "replace into conditions (thing_id, user_id, cond) values (?,?,?)"
	stmtSetCondition, err = u.Sth(db, ssc)
	if err != nil {
		glog.Fatalf("sth(%s): %s", ssc,err)
	}
	sgc := "select IFNULL(cond,0) from conditions where thing_id=? and user_id=?"
	stmtGetCondition, err = u.Sth(db,sgc)
	if err != nil {
		glog.Fatalf("sth(%s): %s", sgc,err)
	}
	sguc := "select thing_id, cond from conditions where user_id=?"
	stmtUserConditions, err = u.Sth(db,sguc)
	if err != nil {
		glog.Fatalf("sth(%s): %s", sguc,err)
	}

}
func (t MyThing) Condition() (c int) {
	err := stmtGetCondition.QueryRow(t.ID,t.Coll.UserID).Scan(&c)
	switch {
		case err == sql.ErrNoRows:
			return 0
		case err != nil:
			glog.Errorf("stmtGetCondition.QueryRow(%v,%v): %s", t.ID, t.Coll.UserID, err)
			return 0
	}
	return c
}
func (t MyThing) SetCondition(c int) (err error) {
	_,err = stmtSetCondition.Exec(t.ID, t.Coll.UserID, c)
	if err != nil {glog.Errorf("stmtSetCondition.Exec(%v,%v,%v): %s",t.ID, t.Coll.UserID, c,err)}
	return err
}
func (coll Collection) GetConditions() (map[int]int) {
	var rm = make(map[int]int)
	rows, err := stmtUserConditions.Query(coll.UserID)
	switch {
		case err == sql.ErrNoRows:
			return rm
		case err != nil:
			glog.Errorf("stmtUserconditions.Query(%v): %s", coll.UserID, err)
			return rm
	}
	for rows.Next(){
		var r int
		var tid int
		rows.Scan(&tid, &r)
		rm[tid]=r
	}
	return rm
}
