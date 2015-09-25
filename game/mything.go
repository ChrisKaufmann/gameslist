package game
// MyThing

import (
	"github.com/golang/glog"
)
const hasbg string = "#c0c0c0"

type MyThing struct {
	Thing
	Coll Collection
}
type PrintableThing struct {
	Thing
	HasManual bool
	Has bool
	HasBox bool
	ManualID int
	BoxID int
}
func (t PrintableThing) Background() (string) {
	if t.Has {return hasbg}
	return "white"
}
func (t PrintableThing) BoxBackground() (string) {
	if t.HasBox {return hasbg}
	return "white"
}
func (t PrintableThing) ManualBackground() (string) {
	if t.HasManual {return hasbg}
	return "white"
}
func (t PrintableThing) Checked() (string) {
	if t.Has {return "checked"}
	return "unchecked"
}
func (t PrintableThing) ManualChecked() (string) {
	if t.HasManual {return "checked"}
	return "unchecked"
}
func (t PrintableThing) BoxChecked() (string) {
	if t.HasBox {return "checked"}
	return "unchecked"
}

// object functions
func (m MyThing) HasManual() (h bool) {
	man := m.Manual()
	return m.Coll.Have(man)
}
func (m MyThing) HasBox() (h bool) {
	box := m.Box()
	return m.Coll.Have(box)
}

// non object functions
func GetPrintableThings(tl []Thing, ph map[int]bool) (ptl []PrintableThing) {
	allthings, err := GetAllThings()
	if err != nil {glog.Errorf("GetPrintableThings-GetAllThings(): %s", err) }
	tms := make(map[int]int)
	tbs := make(map[int]int)
	//make a hash of all manuals and boxes
	for _,t := range allthings {
		if t.Type=="manual" {tms[t.ParentID]=t.ID}
		if t.Type=="box" {tbs[t.ParentID]=t.ID}
	}
	for _, t := range tl {
		pt := PrintableThing{t,false,false,false,tms[t.ID],tbs[t.ID]}
		if ph[t.ID] {pt.Has = true }
		if ph[tms[t.ID]] {pt.HasManual = true }
		if ph[tbs[t.ID]] {pt.HasBox = true }
		ptl = append(ptl, pt)
	}
	return ptl
}
