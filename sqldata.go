package main

import (
	"database/sql"
	"fmt"
	"os"
)

type DatabasePostg struct {
	dbname     string
	table      map[string][]string
	primkey    map[string][]string
	foreignkey map[string][]string
}

func InitConnStr() (dbdriver string, connstr string) {
	// Later we change this init strings, its only for developing
	cstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	dbd := "postgres"
	return dbd, cstr
}

func InitDB(dbdriver string, connstr string) (*sql.DB, error) {
	return sql.Open(dbdriver, connstr)
}

func ReadPgSQLQueue() (m map[int]SQLQueries) {
	mmm := make(map[int]SQLQueries, 3)
	mmm[0] = Tables
	mmm[1] = PrimaryKeys
	mmm[2] = ForeignKeys
	return mmm
}
