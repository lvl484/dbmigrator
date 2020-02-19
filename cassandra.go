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
	//name of keyspace
	clName string
	//map of tables - key:tablename-data - [0] string of create query, [1] string of insert query
	TableInsert map[string][]string
	//active Casssandra cluster
	clust *gocql.ClusterConfig
	//active session
	sess *gocql.Session
}

// CreateCassandra create new instance of Cassandra handle
func CreateCassandra(dbname string) (*Cassandra, error) {
	cluster := gocql.NewCluster(HOST)
	cluster.Consistency = gocql.Quorum
	activesess, err := cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Drop keyspace if exist
	err = activesess.Query(strings.Join([]string{"DROP KEYSPACE IF EXISTS", dbname}, " ")).Exec()
	if err != nil {
		log.Println(err)
	}
	str := strings.Join([]string{"CREATE KEYSPACE", dbname, "WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"}, " ")
	err = activesess.Query(str).Exec()
	if err != nil {
		log.Println(err)
	}
	cluster.Keyspace = dbname

	return &Cassandra{
		clName:      dbname,
		TableInsert: make(map[string][]string),
		sess:        activesess,
		clust:       cluster,
	}, err
}

// CheckSesion check if session is open, if not then Open session, if could't return error
func (cs *Cassandra) CheckSesion() error {
	var err error
	if cs.sess.Closed() {
		activesess, err := cs.clust.CreateSession()
		if err == nil {
			cs.sess = activesess
		}
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
	var tmpcreate, tmpinsdata, tmpinsert []string
	for _, col = range tab.Columns {
		c := ""
		if col.Pk {
			c = "PRIMARY KEY"
		}

		tmpcreate = append(tmpcreate, strings.Join([]string{col.Cname, ConvTypePostgCasan()[col.Ctype], c}, " "))
		tmpinsert = append(tmpinsert, col.Cname)
		tmpinsdata = append(tmpinsdata, "?")
	}
	screate := fmt.Sprintf("CREATE TABLE %s (%s)", tablename, strings.Join(tmpcreate, ","))
	sinsert := fmt.Sprintf("COPY %s (%s) VALUES (%s) FROM STDIN", tablename, strings.Join(tmpinsert, ","), strings.Join(tmpinsdata, ","))
	cs.TableInsert[tablename] = []string{screate, sinsert}
}

// CopyDataToDB write data to Cassandra table "tableName"
func (cs *Cassandra) CopyDataToDB(copyquery string, rows *sql.Rows) error {
	err := cs.CheckSesion()
	if err != nil {
		return err
	}
	if err := cs.sess.Query(copyquery, rows).Exec(); err != nil {
		log.Println(err)
	}
	return err
}

// WriteSchemaToDB runs specific Cassandra query
func (cs *Cassandra) WriteSchemaToDB() error {
	err := cs.CheckSesion()
	if err != nil {
		return err
	}
	for _, tab := range cs.TableInsert {
		if len(tab) > 0 {
			createStr := tab[0]
			if createStr != "" {

				if err := cs.sess.Query(createStr).Exec(); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return err
}
