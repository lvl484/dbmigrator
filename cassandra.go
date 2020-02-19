package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gocql/gocql"
)

const HOST = "127.0.0.1"

// Cassandra structure to handle connection to Cassandra DB
type Cassandra struct {
	//map of tables - key:tablename-data - [0] string of create query, [1] string of insert query
	TableInsert map[string][]string
}

// CreateCassandra create new instance of Cassandra handle
func CreateCassandra() (*Cassandra, error) {
	err := CheckSesion()
	if err != nil {
		return nil, err
	}
	// Drop keyspace if exist
	err = CassandraSession.Query(strings.Join([]string{"DROP KEYSPACE IF EXISTS", DBKeyspace}, " ")).Exec()
	if err != nil {
		log.Println(err)
	}
	str := strings.Join([]string{"CREATE KEYSPACE", DBKeyspace, "WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"}, " ")
	err = CassandraSession.Query(str).Exec()
	if err != nil {
		log.Println(err)
	}
	//	CassandraCluster.Keyspace = DBKeyspace
	return &Cassandra{
		TableInsert: make(map[string][]string),
	}, err
}

// CheckSesion check if session is open, if not then Open session, if could't return error
func CheckSesion() error {
	var err error
	if CassandraSession == nil {
		if CassandraCluster == nil {
			CassandraCluster := gocql.NewCluster(HOST)
			CassandraSession, err = CassandraCluster.CreateSession()
		}
	} else if CassandraSession.Closed() {
		CassandraSession, err = CassandraCluster.CreateSession()
	}
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// CreateTableScheme creates schema of tables in Cassandra
func (cs *Cassandra) CreateTableScheme(dbD *DatabasePostg) error {
	for tablesh, tabc := range dbD.Tables {
		cs.CreateQueryForTable(tablesh, tabc)
	}
	err := cs.WriteSchemaToDB()
	return err
}

// CreateQueryForTable creates string of create and input query for selected table
func (cs *Cassandra) CreateQueryForTable(tablename string, tab Table) {
	var col Column
	var tmpcreate, tmpinsdata, tmpinsert, cp []string
	c := ""
	cf1 := ""
	cf2 := ""
	for _, col = range tab.Columns {

		if col.Pk {
			c = "PRIMARY KEY"
			cp = append(cp, col.Cname)
			cf1 = "("
			cf2 = ")"
		}

		tmpcreate = append(tmpcreate, strings.Join([]string{col.Cname, ConvTypePostgCasan()[col.Ctype], c}, " "))
		tmpinsert = append(tmpinsert, col.Cname)
		tmpinsdata = append(tmpinsdata, "?")
	}
	cpi := strings.Join(cp, ",")
	tmpprimary := strings.Join([]string{c, cf1, cpi, cf2}, " ")
	tmpcreate = append(tmpcreate, tmpprimary)
	screate := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (%s)", DBKeyspace, tablename, strings.Join(tmpcreate, ","))
	sinsert := fmt.Sprintf("COPY %s.%s (%s) VALUES (%s) FROM STDIN", DBKeyspace, tablename, strings.Join(tmpinsert, ","), strings.Join(tmpinsdata, ","))
	cs.TableInsert[tablename] = []string{screate, sinsert}
}

// CopyDataToDB write data to Cassandra table "tableName"
func (cs *Cassandra) CopyDataToDB(copyquery string, rows *sql.Rows) error {
	err := CheckSesion()
	if err != nil {
		return err
	}
	if err := CassandraSession.Query(copyquery, rows).Exec(); err != nil {
		log.Println(err)
	}
	return err
}

// WriteSchemaToDB runs specific Cassandra query
func (cs *Cassandra) WriteSchemaToDB() error {
	err := CheckSesion()
	if err != nil {
		return err
	}
	for _, tab := range cs.TableInsert {
		if len(tab) > 0 {
			createStr := tab[0]
			if createStr != "" {

				if err := CassandraSession.Query(createStr).Exec(); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return err
}
