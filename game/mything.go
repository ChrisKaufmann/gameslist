package game
// MyThing

import (
)

type MyThing struct {
	Thing
	Coll Collection
}
type PrintMyThing struct {
	MyThing
	Background string
	ManualBackground string
	BoxBackground string
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
func tf(i int)(b bool) {
	if i > 0 { return true }
	return false
}
