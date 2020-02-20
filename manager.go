package main

import (
	"sync"
)

// Manager manage a proces of migration SQL entities to NOSQL databases and vice versa
type Manager struct {
	posg *SQLPostgres
	cass *Cassandra
	wg   sync.WaitGroup
}

// NewManager create a new instance of structure Manager
func NewManager() (*Manager, error) {
	pgdb, err := newSQLPostgres()
	if err != nil {
		return nil, err
	}
	cassand, err := NewCassandra()
	if err != nil {
		return nil, err
	}
	return &Manager{
		cass: cassand,
		posg: pgdb,
	}, err
}

// SetKeyspace set a keyspace and name of DB to our structure
func (mn *Manager) SetKeyspace(keyspasename string) {
	mn.posg.DbData.Databasename = keyspasename
	mn.cass.DBKeyspace = keyspasename
}

// GetSchemaFromSQl takes schema of tables from SQL entries
func (mn *Manager) GetSchemaFromSQL() error {
	mn.wg.Add(1)
	keyspname, err := mn.posg.ReadDBName()
	if err != nil {
		return err
	}
	mn.SetKeyspace(keyspname)
	err = mn.posg.ReadTableSchema()
	if err != nil {
		return err
	}
	err = mn.posg.ReadPrimaryKeys()
	if err != nil {
		return err
	}

	err = mn.posg.ReadForeignKeys()
	if err != nil {
		return err
	}
	defer mn.wg.Done()
	return err
}

// PutSchemaToNoSQL create structure in NoSQL entry according to schems from SQL entries
func (mn *Manager) PutSchemaToNoSQL() error {
	var err error
	mn.wg.Add(1)
	err = mn.cass.CreateKeyspaceCassandra()
	if err != nil {
		return err
	}
	err = mn.cass.CreateTableScheme(mn.posg.DbData)
	if err != nil {
		return err
	}
	defer mn.wg.Done()
	return err
}

// GetDataFromSQL reading data from SQL entries according to schema
func (mn *Manager) GetDataFromSQL() error {
	var err error
	for tablename, tab := range mn.cass.TableQueries {
		errchanel := make(chan error, 1)
		go func() {
			errchanel <- mn.ReaDataFromSingleTable(tablename, tab.QueryInsert)

		}()
		err = <-errchanel
		if err != nil {
			return err
		}
	}
	mn.wg.Wait()
	return err
}

// ReaDataFromSingleTable read data from SQL table and write it to NoSQL
func (mn *Manager) ReaDataFromSingleTable(tablename string, tableInsQuery string) error {
	mn.wg.Add(1)
	rows, err := mn.posg.ReadDataFromTable(tablename)
	if err != nil {
		return err
	}
	err = mn.cass.CassandraSession.Query(tableInsQuery, rows).Exec()
	if err != nil {
		return err
	}
	mn.wg.Done()
	return err
}

func (mn *Manager) CloseConnection() {
	mn.cass.CassandraSession.Close()
	mn.posg.pdb.Close()
}
