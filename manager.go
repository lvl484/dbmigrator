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
	return &Manager{
		posg: pgdb,
	}, err
}

// GetSchemaFromSQl takes schema of tables from SQL entries
func (mn *Manager) GetSchemaFromSQl() error {
	mn.wg.Add(1)
	err := mn.posg.ReadDBName()
	if err != nil {
		return err
	}
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
	mn.cass, err = CreateCassandra()
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
	mn.wg.Add(1)
	go func() error {
		var err error
		for tablename, tsl := range mn.cass.TableInsert {
			rows, err := mn.posg.ReadDataFromTable(tablename)
			if err != nil {
				return err
			}
			err = mn.cass.CopyDataToDB(tsl[1], rows)
			if err != nil {
				return err
			}
		}
		return err
	}()
	defer mn.wg.Done() //??????????
	mn.wg.Wait()
	return err
}
