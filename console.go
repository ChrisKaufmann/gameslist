package main
//Consoles - crud for consoles

import (
  "database/sql"
)
var (
	stmtGetConsoleList		*sql.Stmt
	stmtSaveConsole			*sql.Stmt
)
type Console struct {
	ID		int
	Name	string
}
func (c  Console) Print() {
	print("ID:\t"+tostr(c.ID)+"\nName:\t"+c.Name+"\n")
}
func (c Console) Save(){
	if c.ID > 0 {
		stmtSaveConsole.Exec(c.Name)
	} else {
		_, err := stmtAddConsole(c.Name)
		if err != nil {
			err.Error()
		}
	}
}
func allConsoles() []Console {
	var cl []Console
	var stmt =sth("select id, name from console where 1")
	rows, err := stmt.Query()
	if err != nil {
		err.Error()
		return cl
	}
	for rows.Next() {
		var c Console
		rows.Scan(&c.ID,&c.Name)
		cl = append(cl,c)
	}
	return cl
}
func getConsole(id string)  Console {
	if id == "" {
		return c
	}
	for _,c := range allConsoles() {
		if tostr(c.ID) == id {
			return c
		}
	}
	return c
}
