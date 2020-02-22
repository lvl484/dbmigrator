package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

const HOST = "127.0.0.1"
const cCBootstrapTimeout time.Duration = 9 * time.Second

// TableQuery structure that holds two CQL queries for creating table and inserting data to it
type TableQuery struct {
	QuerySelect string
	QueryCreate string
	QueryInsert string
}

// Cassandra structure to handle connection to Cassandra DB
type Cassandra struct {
	CassandraSession *gocql.Session
	DBKeyspace       string
	CassandraCluster *gocql.ClusterConfig
	//map of tables - with query
	TableQueries map[string]TableQuery
}

// NewCassandra create new instance that manage session with db Cassandra
func NewCassandra() (*Cassandra, error) {
	cluster := gocql.NewCluster(HOST)
	cluster.Timeout = cCBootstrapTimeout
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &Cassandra{
		CassandraSession: session,
		CassandraCluster: cluster,
		TableQueries:     make(map[string]TableQuery),
	}, nil
}

// CreateKeyspaceCassandra create a keyspace of Cassandra handle
func (cs *Cassandra) CreateKeyspaceCassandra() error {
	str := fmt.Sprintf(KeyspaceQuery, cs.DBKeyspace)
	err := cs.CassandraSession.Query(str).Exec()
	if err != nil {
		return err
	}
	return nil
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
	var tmpselect, tmpcreate, tmpinsdata, tmpinsert, cp []string
	for _, col = range tab.Columns {
		if col.Pk {
			cp = append(cp, col.Cname)
		}
		selectcolumn := "%s"
		if CompareType(col.Ctype) {
			selectcolumn = "CAST (%s as real)"
		}
		tmpselect = append(tmpselect, fmt.Sprintf(selectcolumn, col.Cname))
		tmpcreate = append(tmpcreate, strings.Join([]string{col.Cname, ConvTypePostgCasan()[col.Ctype]}, " "))
		tmpinsert = append(tmpinsert, col.Cname)
		tmpinsdata = append(tmpinsdata, "?")
	}
	tmpprimary := fmt.Sprintf(CassandraPrimary, strings.Join(cp, ","))
	queryselect := fmt.Sprintf(PostgresSelect, strings.Join(tmpselect, ","), tablename)
	tmpcreate = append(tmpcreate, tmpprimary)
	screate := fmt.Sprintf(CassandraTable, cs.DBKeyspace, tablename, strings.Join(tmpcreate, ","))
	sinsert := fmt.Sprintf(CassandraCopyData, cs.DBKeyspace, tablename, strings.Join(tmpinsert, ","), strings.Join(tmpinsdata, ","))
	cs.TableQueries[tablename] = TableQuery{queryselect, screate, sinsert}
}

// WriteSchemaToDB runs specific Cassandra query
func (cs *Cassandra) WriteSchemaToDB() error {
	var err error
	for _, tab := range cs.TableQueries {
		createStr := tab.QueryCreate
		if createStr != "" {
			err = cs.CassandraSession.Query(createStr).Exec()
			if err != nil {
				return err
			}
		}
	}
	return err
}
