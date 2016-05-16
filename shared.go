package main

import (
	"./game"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"time"
)

var (
	ShareMain             = template.Must(template.ParseFiles("templates/ShareMain.html"))
	ReadTableEntryConsole = template.Must(template.ParseFiles("templates/ReadTableEntryConsole.html"))
	ReadTableEntryGame    = template.Must(template.ParseFiles("templates/ReadTableEntryGame.html"))
)

func handleShared(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	var sid string
	u.PathVars(r, "/share/", &sid)

	suser, err := auth.GetUserByShared(sid)
	if err != nil {
		glog.Errorf("handleShared - auth.GetUserByShared(%s): %s", sid, err)
		return
	}

	cons, err := game.GetConsoles(suser)
	if err != nil {
		glog.Errorf("handleShared - game.GetAllConsoles(): %s", err)
		return
	}

	m := make(map[game.Console][]game.Game)
	for _, c := range cons {
		gl, err := c.Games()
		if err != nil {
			glog.Errorf("c.Games(): %s", err)
			return
		}
		for _, g := range gl {
			if g.Has {
				m[c] = append(m[c], g)
			}
		}
	}
	fmt.Printf("handleShared %v\n", time.Now().Sub(t0))
}
