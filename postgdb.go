package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type SQLPostgres struct {
	pdb    *sql.DB
	DbData *DatabasePostg
}

// newSQLPostgres create a new instance of SQLPostgres which provide connection with Postgresql DB
func newSQLPostgres() (*SQLPostgres, error) {
	cstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s ",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	dbd := "postgres"
	database, err := sql.Open(dbd, cstr)

	if err != nil {
		return nil, err
	}
	datapostg := NewDatabasePostg()
	return &SQLPostgres{
		pdb:    database,
		DbData: datapostg,
	}, nil
}
func (sp *SQLPostgres) DatabaseSQL() (*sql.DB, error) {
	err := sp.pdb.Ping()
	if err != nil {
		return nil, err
	}
	return sp.pdb, err
}

// AddPrimaryKey add PRIMARY KEY to table by column name
func (sp *SQLPostgres) AddPrimaryKey(tabl string, col string) {
	t := sp.DbData.Tables[tabl]
	colum := t.Columns
	if len(colum) == 0 {
		log.Printf("Failed to create PRIMARY KEY on Table: %s column: %s", tabl, col)
	}
	for i, c := range colum {
		if c.Cname == col {
			c.MakePrimaryKey()
			colum[i] = c
			sp.DbData.Tables[tabl] = t
			break
		}
	}

}

// AddForeignKey add Foreign keys by it's constraint name
func (sp *SQLPostgres) AddForeignKey(constrname string, forkey ForeignKey) {
	sp.DbData.Foreignkeys[constrname] = forkey
}

// AddColumnToTable add column to table
func (sp *SQLPostgres) AddColumnToTable(tname string, col Column) {
	tab := sp.DbData.Tables[tname]
	tab.AddColumn(col)
	sp.DbData.Tables[tname] = tab
}

// ReadDBName read DB name
func (sp *SQLPostgres) ReadDBName() error {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return err
	}
	var ss string
	err = db.QueryRow(string(DbNameQuery)).Scan(&ss)
	if err != nil {
		return err
	}
	sp.DbData.SetDBName(ss)

	return err
}

// ReadTableSchema reads schema of tables from SQL entry
func (sp *SQLPostgres) ReadTableSchema() error {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return err
	}
	rows, err := db.Query(TablesQuery, sp.DbData.Databasename)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tn, cn, dt string
		var pos int
		err = rows.Scan(&tn, &pos, &cn, &dt)
		if err != nil {
			log.Println(err)
		}
		col := Column{Cname: cn, Ctype: dt, Pk: false}
		sp.AddColumnToTable(tn, col)
	}
	return err
}

// ReadPrimaryKeys reads Primary keys of tables from SQL entry
func (sp *SQLPostgres) ReadPrimaryKeys() error {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return err
	}
	rows, err := db.Query(PrimaryKeysQuery, sp.DbData.Databasename)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var st, sc string
		err = rows.Scan(&st, &sc)
		if err != nil {
			return err
		}
		sp.AddPrimaryKey(st, sc)
	}
	return err
}

// ReadForeignKeys reads foreign key
func (sp *SQLPostgres) ReadForeignKeys() error {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return err
	}
	rows, err := db.Query(ForeignKeysQuery, sp.DbData.Databasename)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var f ForeignKey
		var s string
		err = rows.Scan(&s, &f.Tablename, &f.Colname, &f.Tableref, &f.Colref)
		if err != nil {
			return err
		}
		sp.AddForeignKey(s, f)
	}
	return err
}

// ReadForeignKeys reads foreign key
func (sp *SQLPostgres) ReadDataFromTable(tablename string) (*sql.Rows, error) {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(DataTablesQuery, tablename)
	fmt.Printf("Reading data from %s", tablename) // it's for testing
	if err != nil {
		return nil, err
	}
	return rows, err
}
