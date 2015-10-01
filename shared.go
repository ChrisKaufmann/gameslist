package main

import (
	"./auth"
	"fmt"
	"./game"
	"time"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"net/http"
)
var (
	ReadTableEntryConsole = template.Must(template.ParseFiles("templates/ReadTableEntryConsole.html"))
	ReadTableEntryGame = template.Must(template.ParseFiles("templates/ReadTableEntryGame.html"))
)

func handleShared(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	var sid string
	u.PathVars(r,"/shared/", &sid)

	suser, err := auth.GetUserByShared(sid)
	if err != nil { glog.Errorf("handleShared - auth.GetUserByShared(%s): %s", sid, err);return }

	coll, err := game.GetCollection(suser.ID)
	if err != nil { glog.Errorf("handleShared - game.GetCollection(%s): %s", suser.ID, err);return }

	cons, err := game.GetAllConsoles()
	if err != nil { glog.Errorf("handleShared - game.GetAllConsoles(): %s", err); return }

	cl, err := coll.Things()
	if err != nil {glog.Errorf("handleShared - coll.Things(): %s", err); return}

	mtl := coll.MyThingsHash()
	mttl := game.GetPrintableThings(cl, mtl)

	fmt.Fprintf(w,"<table>")
	fmt.Fprintf(w,"<tr><td colspan=2><a name='sym'></a>Console</td><td align=right>Game</td><td>?</td><td>Man</td><td>Box</td></tr>")
	for _, myc := range game.GetPrintableThings(cons, mtl) {
		ReadTableEntryConsole.Execute(w,myc)
		for _,t := range mttl {
			if t.ParentID == myc.ID {
				ReadTableEntryGame.Execute(w,t)
			}
		}
	}
	fmt.Fprintf(w,"</table>")
	fmt.Printf("handleShared %v\n", time.Now().Sub(t0))
}
