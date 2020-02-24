package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// SQLPostgres structure that manage connection to PosgreSQL
type SQLPostgres struct {
	pdb    *sql.DB
	DbData *DatabasePostg
}

// newSQLPostgres create a new instance of SQLPostgres which provide connection with Postgresql DB
func newSQLPostgres() (*SQLPostgres, error) {

	cstr := fmt.Sprintf(PostgresConfDB, os.Getenv("MIGRATION_HOST"), os.Getenv("MIGRATION_PORT"), os.Getenv("MIGRATION_USER"), os.Getenv("MIGRATION_PASSWORD"), os.Getenv("MIGRATION_DBNAME"))
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

// DatabaseSQL check connection with DB
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
	//	fmt.Printf("%s  %s   %s\n", tname, col.Cname, col.Ctype)
}

// ReadDBName read DB name
func (sp *SQLPostgres) ReadDBName() (string, error) {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return "", err
	}
	var ss string
	err = db.QueryRow(string(DbNameQuery)).Scan(&ss)
	if err != nil {
		return "", err
	}
	return ss, err
}

// ReadTableSchema reads schema of tables from SQL entry
func (sp *SQLPostgres) ReadTableSchema() error {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return err
	}
	rows, err := db.Query(TablesSchemaQuery, sp.DbData.Databasename)
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
		var tablename, columnname string
		err = rows.Scan(&tablename, &columnname)
		if err != nil {
			return err
		}
		sp.AddPrimaryKey(tablename, columnname)
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
		var foreignkey ForeignKey
		var constraintname string
		err = rows.Scan(&constraintname, &foreignkey.Tablename, &foreignkey.Colname, &foreignkey.Tableref, &foreignkey.Colref)
		if err != nil {
			return err
		}
		sp.AddForeignKey(constraintname, foreignkey)
	}
	return err
}

// ReadForeignKeys reads foreign key
func (sp *SQLPostgres) ReadDataFromTable(querystring string) (*sql.Rows, error) {
	db, err := sp.DatabaseSQL()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(querystring)
	if err != nil {
		return nil, err
	}
	return rows, err
}

// CreateNewDatabase creates new database for writing data in
func (sp *SQLPostgres) CreateNewDatabase(dbname string) error {
	var err error
	cstr := "user=postgres host=localhost password=root sslmode=disable"
	database, _ := sql.Open("postgres", cstr)
	dbname += "v1"

	creatstr := fmt.Sprintf(`create database %s`, dbname)
	_, err = database.Exec(creatstr)

	if err.(*pq.Error).Code == "42P04" {
		strquery := "DROP SCHEMA public CASCADE"
		_, err = database.Exec(strquery)
		if err != nil {
			return err
		}
		strquery = "CREATE SCHEMA public"
		_, err = database.Exec(strquery)
		if err != nil {
			return err
		}
		strquery = "GRANT ALL ON SCHEMA public TO postgres"
		_, err = database.Exec(strquery)
		if err != nil {
			return err
		}
		strquery = "GRANT ALL ON SCHEMA public TO public"
		_, err = database.Exec(strquery)
		if err != nil {
			return err
		}

		err = nil
	}
	sp.pdb = database
	return err
}

// WriteSchemaToDatabase creates tabels in SQL entry according schema from NoSQL
func (sp *SQLPostgres) WriteSchemaToDatabase(tabs map[string]TableQuery) error {
	db, err := sp.DatabaseSQL()
	for _, tab := range tabs {
		_, err = db.Exec(tab.QueryCreate)
		if err != nil {
			return err
		}
	}
	return err
}
