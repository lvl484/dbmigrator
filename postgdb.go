package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DatabasePostg struct {
	dbname         string
	chset          string
	table          map[string][]string
	dbfunc         map[string][]string
	dbfuncparam    map[string][]string
	sequence       map[string][]string
	primkey        map[string][]string
	foreignkey     map[string][]string
	checkconst     map[string][]string
	view           map[string][]string
	autoincrements map[string][]int
	triger         map[string][]string
}

type SQLQuery struct {
	Order  int    `json:"id"`
	StrSQL string `json:"SQL"`
}

type SQLPostgres struct {
	pdb    *sql.DB
	dbData *DatabasePostg
	qrSQL  map[int]string //Queries
}

func NewSQLPostgre() *SQLPostgres {
	dbdrv, constr := InitConnStr()

	return &SQLPostgres{
		pdb:   InitDB(dbdrv, constr),
		qrSQL: ReadPgSQLQueue("pgread.json"),
	}
}

func InitConnStr() (dbdriver string, connstr string) {
	// Later we change this init strings, its only for developing
	cstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "root", "dvdrental")
	dbd := "postgres"
	return dbd, cstr
}

func InitDB(dbdriver string, connstr string) *sql.DB {
	database, err := sql.Open(dbdriver, connstr)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func ReadPgSQLQueue(fj string) (m map[int]string) {

	jsFile, err := os.Open(fj)
	if err != nil {
		log.Print(err)
	}
	defer jsFile.Close()

	byteValue, _ := ioutil.ReadAll(jsFile)
	var mm map[string]*json.RawMessage

	err = json.Unmarshal(byteValue, &mm)
	if err != nil {
		log.Print(err)
	}
	var queries []SQLQuery
	err = json.Unmarshal(*mm["SQLQuery"], &queries)
	if err != nil {
		log.Print(err)
	}
	mmm := make(map[int]string)
	for _, q := range queries {
		mmm[q.Order] = q.StrSQL
	}
	return mmm
}

func (sp *SQLPostgres) GetDataFromSource() {
	var strQuery string

	for i := 0; i < 15; i++ {
		strQuery := sp.qrSQL[i]
		go sp.ExecQuery(i, strQuery)

	}

}

func (sp *SQLPostgres) ExecQuery(ord int, sq string) {
	ch := make(chan string)
	defer close(ch)

	qarg := sp.selectArg(ord)

	rows, err := sp.pdb.Query(sq, qarg)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan()

		if err != nil {
			log.Println(err)
			continue
		}

	}
}

func (sp *SQLPostgres) selectArg(ord int) []string {
	var gArg []string
	switch ord {
	case 1:
		gArg = nil
	case 3:
		gArg = nil
	case 7:
	case 8:
	case 9:
	case 10:
	case 11:
	case 12:
	case 13:
	case 14:
	default:
		gArg := [1]string{sp.dbData.dbname}

	}
	return gArg
}

/*
{
    "_comment": "list all tables",
    "id": 7,
    "SQL": "SELECT *FROM information_schema.tables;"
},
{
    "_comment": "lists all views in the database",
    "id": 8,
    "SQL": "SELECT * FROM information_schema.views;"
},
{
    "_comment": "lists all constraints from tables in this database",
    "id": 10,
    "SQL": "SELECT * FROM information_schema.table_constraints;"
},
{
    "_comment": "lists all foreign keys in the database",
    "id": 11,
    "SQL": "SELECT * FROM information_schema.referential_constraints;"
},
{
    "_comment": "Cataloglist all check constraints",
    "id": 12,
    "SQL": "SELECT * FROM information_schema.check_constraints;"
},
{
    "_comment": "list all triggers",
    "id": 14,
    "SQL": "SELECT *FROM information_schema.triggers;"
},*/
