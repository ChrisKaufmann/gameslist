package main

import (
	"./game"
	"fmt"
	"github.com/ChrisKaufmann/goauth"
	u "github.com/ChrisKaufmann/goutils"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"sort"
	"time"
)

var (
	sharedtmpl = template.Must(template.ParseGlob("templates/shared.html"))
)

func handleShared(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	var sid string
	u.PathVars(r, "/share/", &sid)

	user, err := auth.GetUserByShared(sid)
	if err != nil {
		glog.Errorf("handleShared - auth.GetUserByShared(%s): %s", sid, err)
		return
	}
	cml := make(map[string][]game.Game)
	gl := game.UserGames(user)
	sort.Sort(game.GameName(gl))
	for _, g := range gl {
		cml[g.ConsoleName] = append(cml[g.ConsoleName], g)
	}
	cl, err := game.GetConsoles(user)
	var sm []game.ConsoleMeta
	for _, c := range cl {
		if len(cml[c.Name]) > 0 {
			var cm game.ConsoleMeta
			cm.Console = c
			cm.Games = cml[c.Name]
			sm = append(sm, cm)
		}
	}
	sort.Sort(game.ConsoleMetaName(sm))

	sharedtmpl.ExecuteTemplate(w, "shared", sm)

	fmt.Printf("handleShared %v\n", time.Now().Sub(t0))
}
