package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type SQLPostgres struct {
	pdb    *sql.DB
	cass   *Cassandra
	dbData *DatabasePostg
	qrSQL  map[int]SQLQueries
}

func NewSQLPostgre() *SQLPostgres {
	dbdrv, constr := InitConnStr()
	db, err := InitDB(dbdrv, constr)
	if err != nil {
		log.Println(err)
	}
	return &SQLPostgres{
		pdb:   db,
		qrSQL: ReadPgSQLQueue(),
	}
}

func (sp *SQLPostgres) GetSchemaFromSQl() {
	sp.ExecPGQuery(0, DbName)
	for i, q := range sp.qrSQL {
		go sp.ExecPGQuery(i+1, q)
	}
}

func (sp *SQLPostgres) PutSchemaToNoSQL() {
	for _, tabs := range sp.dbData.table {
		i, ss := sp.CreatePrimaryField(tabs[0])
		ss = append(tabs, ss...)
		go sp.WriteToCassandra(1, i, ss)
	}
	go sp.WriteToCassandra(2, 0, nil)

}

func (sp *SQLPostgres) CreatePrimaryField(tabl string) (int, []string) {
	var ss []string
	for _, pk := range sp.dbData.primkey[tabl] {
		ss = append(ss, fmt.Sprintf("PRIMARY KEY ( %s )", pk))
	}
	i := len(sp.dbData.primkey[tabl])
	return i, ss
}

func (sp *SQLPostgres) GetDataFromSQL() {

	go sp.ExecPGQuery(3, DataTables)
}

func (sp *SQLPostgres) PutDataToNoSQL(tn string, rows *sql.Rows) {
	ch := make(chan *sql.Rows)
	defer close(ch)
	ch <- rows
	go sp.cass.InputData(tn, ch)
}

func (sp *SQLPostgres) ExecPGQuery(order int, sq SQLQueries) {
	switch order {
	case 0:
		err := sp.pdb.QueryRow(string(sq)).Scan(sp.dbData.dbname)
		if err != nil {
			log.Println(err)
		}
		sp.cass = CreateCassandra(sp.dbData.dbname)
	case 3:
		for _, tab := range sp.dbData.table {
			tn := tab[0]
			rows, err := sp.pdb.Query(string(sq), tn)
			if err != nil {
				log.Println(err)
			}
			defer rows.Close()
			sp.PutDataToNoSQL(tn, rows)
		}

	default:
		dbnam := sp.dbData.dbname
		rows, err := sp.pdb.Query(string(sq), dbnam)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		s := make([]string, 4)
		ch := make(chan []string, 1)
		defer close(ch)

		for rows.Next() {
			err = rows.Scan(s)
			if err != nil {
				log.Println(err)
			}
			ch <- s
			go sp.CreateData(order, ch)
		}
	}
}

func (sp *SQLPostgres) CreateData(order int, chOut <-chan []string) {
	ss := <-chOut
	switch order {
	case 1:
		sp.dbData.table[ss[0]] = append(sp.dbData.table[ss[0]], ss[2:]...)
	case 2:
		sp.dbData.primkey[ss[0]] = append(sp.dbData.primkey[ss[0]], ss[1:]...)
	case 3:
		sp.dbData.foreignkey[ss[0]] = append(sp.dbData.foreignkey[ss[0]], ss[1:]...)
	}

}

func (sp *SQLPostgres) WriteToCassandra(order int, parnum int, ds []string) {
	ch := make(chan []string, 1)
	defer close(ch)
	ch <- ds
	go sp.cass.RunQuery(order, parnum, ch)
}
