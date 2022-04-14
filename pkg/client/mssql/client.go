package mssql

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func NewMSSQL() (db *sql.DB, err error) {
	//db, err := sql.Open("mssql",databaseURL)
	db, err = sql.Open("mssql", "sqlserver://cim:cim@10.1.14.21/Panacim?database=PanaCIM&encrypt=disable")
	if err != nil {
		//return nil, err
		log.Fatal("Error creating connerction pool: " + err.Error())
		return nil, err
	}
	//defer db.Close()
	//log.Printf("Connected!\n")
	return //db, nil
}
