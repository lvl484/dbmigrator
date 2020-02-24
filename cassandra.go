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

// CreateNoSQLTableScheme creates schema of tables in Cassandra
func (cs *Cassandra) CreateNoSQLTableScheme(dbD *DatabasePostg) {
	for tablesh, tabc := range dbD.Tables {
		cs.CreateQueryForNoSQLTable(tablesh, tabc)
	}
}

// CreateQueryForSQLTable creates string of create and input query for selected table in Cassandra
func (cs *Cassandra) CreateQueryForSQLTable(tablename string, tab Table) {
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

//ReadTableSchema reads schems of tables in existing keyspace
func (cs *Cassandra) ReadTableSchema(keyspname string, datapos *DatabasePostg) error {
	querystr := fmt.Sprintf(CassandraTableSchema, keyspname)
	iter := cs.CassandraSession.Query(querystr).Iter()
	var tablename string
	for iter.Scan(&tablename) {
		err := cs.ReadColumnSchema(keyspname, tablename, datapos)
		if err != nil {
			return err
		}
	}
	err := iter.Close()
	if err != nil {
		return err
	}
	cs.CreateNoSQLTableScheme(datapos)
	return err
}

// ReadColumnSchema reada column schema from NoSQL entry
func (cs *Cassandra) ReadColumnSchema(keyspname string, tablename string, datapos *DatabasePostg) error {
	querystr := fmt.Sprintf(CassandraColumnSchema, keyspname, tablename)
	iter := cs.CassandraSession.Query(querystr).Iter()
	var casscolumn Column
	var colkind string
	tab := datapos.Tables[tablename]
	for iter.Scan(&casscolumn.Cname, &casscolumn.Ctype, &colkind) {
		casscolumn.Pk = false
		if colkind == "partition_key" || colkind == "clustering" {
			casscolumn.MakePrimaryKey()
		}
		tab.AddColumn(casscolumn)
	}
	datapos.Tables[tablename] = tab
	err := iter.Close()
	return err
}

// CreateSQLTableScheme creates schema of tables in PostgreSQL
func (cs *Cassandra) CreateSQLTableScheme(dbD *DatabasePostg) error {
	for tablename, tabc := range dbD.Tables {
		cs.CreateQueryForSQLTable(tablename, tabc)
	}
	err := cs.WriteSchemaToDB()
	return err
}

// CreateQueryForNoSQLTable creates string of create and input query for selected table in PostgreSQL
func (cs *Cassandra) CreateQueryForNoSQLTable(tablename string, tab Table) {
	var tmpselect, tmpcreate, tmpinsdata, tmpinsert, cp []string
	tablenamesp := fmt.Sprintf("%s.%s", cs.DBKeyspace, tablename)
	newtablenamesp := fmt.Sprintf("%s.%s", "public", tablename)
	for i, col := range tab.Columns {
		if col.Pk {
			cp = append(cp, col.Cname)
		}
		tmpselect = append(tmpselect, col.Cname)
		tmpcreate = append(tmpcreate, strings.Join([]string{col.Cname, ConvTypeCasanPost()[col.Ctype]}, " "))
		tmpinsert = append(tmpinsert, col.Cname)
		tmpinsdata = append(tmpinsdata, fmt.Sprintf("$%d", i+1))
	}
	tmpprimary := fmt.Sprintf(CassandraPrimary, strings.Join(cp, ","))
	tmpcreate = append(tmpcreate, tmpprimary)
	queryselect := fmt.Sprintf(PostgresSelect, strings.Join(tmpselect, ","), tablenamesp)
	screate := fmt.Sprintf(PostgreSQLTable, newtablenamesp, strings.Join(tmpcreate, ","))
	sinsert := fmt.Sprintf(PostgreSQLInsertData, newtablenamesp, strings.Join(tmpinsert, ","), strings.Join(tmpinsdata, ","))
	cs.TableQueries[tablename] = TableQuery{queryselect, screate, sinsert}
}
