package game
// MyThing

import (
	"github.com/golang/glog"
)
const hasbg string = "#c0c0c0"

type MyThing struct {
	Thing
	Coll Collection
	HasManual bool
	Has bool
	HasBox bool
	ManualID int
	BoxID int
	rating int
	review string
}
func (t MyThing) Background() (string) {
	if t.Has {return hasbg}
	return "white"
}
func (t MyThing) BoxBackground() (string) {
	if t.HasBox {return hasbg}
	return "white"
}
func (t MyThing) ManualBackground() (string) {
	if t.HasManual {return hasbg}
	return "white"
}
func (t MyThing) Checked() (string) {
	if t.Has {return "checked"}
	return "unchecked"
}
func (t MyThing) ManualChecked() (string) {
	if t.HasManual {return "checked"}
	return "unchecked"
}
func (t MyThing) BoxChecked() (string) {
	if t.HasBox {return "checked"}
	return "unchecked"
}
func (t MyThing) ConsoleTotal() (int) {
	gl, err := t.Games()
	if err != nil {glog.Errorf("t.Games(): %s", err) }
	return len(gl)
}
func (t MyThing) ConsoleOwned() (int) {
	gl, err := t.Coll.ConsoleGames(t.Thing)
	if err != nil {glog.Errorf("coll.ConsoleGames(%s): %s", t.Thing, err) }
	return len(gl)
}

// non object functions


//for sorting
type ByName []MyThing
func (a ByName) Len() int			{return len(a)}
func (a ByName) Swap(i, j int)		{a[i], a[j] = a[j],a[i]}
func (a ByName) Less(i, j int) bool	{return a[i].Name < a[j].Name }
