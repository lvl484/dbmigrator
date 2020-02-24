package main

import (
	"fmt"
	"sync"
)

// Manager manage a proces of migration SQL entities to NOSQL databases and vice versa
type Manager struct {
	posg *SQLPostgres
	cass *Cassandra
	wg   sync.WaitGroup
}

// NewManager create a new instance of structure Manager
func NewManager(ver int) (*Manager, error) {
	var pgdb *SQLPostgres
	var err error
	switch ver {
	case 1:
		pgdb, err = newSQLPostgres()
		if err != nil {
			return nil, err
		}
	case 2:
		pgdb = &SQLPostgres{DbData: NewDatabasePostg()}
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
	err = mn.cass.CreateSQLTableScheme(mn.posg.DbData)
	if err != nil {
		return err
	}
	defer mn.wg.Done()
	return err
}

// GetDataFromSQL reading data from SQL entries according to schema
func (mn *Manager) GetDataFromSQL() error {
	var err error
	for _, tab := range mn.cass.TableQueries {
		errchanel := make(chan error, 1)
		go func() {
			errchanel <- mn.ReaDataFromSingleTable(tab.QuerySelect, tab.QueryInsert)
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
func (mn *Manager) ReaDataFromSingleTable(selectquery string, insertquery string) error {
	mn.wg.Add(1)
	rows, err := mn.posg.ReadDataFromTable(selectquery)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	pointers := make([]interface{}, len(cols))
	container := make([]interface{}, len(cols))
	for i := 0; i < len(pointers); i++ {
		pointers[i] = &container[i]
	}
	for rows.Next() {
		err = rows.Scan(pointers...)
		if err != nil {
			return err
		}
		err = mn.cass.CassandraSession.Query(insertquery, container...).Exec()
		if err != nil {
			return err
		}
	}
	mn.wg.Done()
	return err
}

// CloseConnection closes all connections to our databases
func (mn *Manager) CloseConnection() {
	mn.cass.CassandraSession.Close()
	err := mn.posg.pdb.Close()
	if err != nil {
		fmt.Println(err)
	}

}

// DoMigrateFromCassandra mannage proces of migration from Cassandra to PostgreSQL
func DoMigrateFromPostgres() error {
	var MigratorManager *Manager
	MigratorManager, err := NewManager(1)
	if err != nil {
		fmt.Printf("There is some problems with inicialization: %v\n", err)
		return err
	}
	defer MigratorManager.CloseConnection()

	err = MigratorManager.GetSchemaFromSQL()
	if err != nil {
		fmt.Printf("Can not get shchema from PostgreSQL: %v\n", err)
		return err
	}

	err = MigratorManager.PutSchemaToNoSQL()
	if err != nil {
		fmt.Printf("Can not put shchema to Cassandra: %v\n", err)
		return err
	}
	err = MigratorManager.GetDataFromSQL()
	if err != nil {
		fmt.Printf("There is some problem with getting data from database: %v\n", err)
		return err
	}
	return err
}

// DoMigrateFromCassandra mannage proces of migration from Cassandra to PostgreSQL
func DoMigrateFromCassandra(cassName string) error {

	var MigratorManager *Manager
	MigratorManager, err := NewManager(2)
	if err != nil {
		fmt.Printf("There is some problems with inicialization: %v\n", err)
		return err
	}

	err = MigratorManager.GetSchemaFromNoSQL(cassName)
	if err != nil {
		return err
	}
	err = MigratorManager.posg.CreateNewDatabase(MigratorManager.cass.DBKeyspace)
	if err != nil {
		return err
	}

	err = MigratorManager.posg.WriteSchemaToDatabase(MigratorManager.cass.TableQueries)
	if err != nil {
		return err
	}

	err = MigratorManager.ManageDataNoSQLtoSQL(MigratorManager.posg)
	if err != nil {
		return err
	}

	defer MigratorManager.CloseConnection()
	return err

}

// ManageDataNoSQLtoSQL start transfer data from NoSQL to SQL database
func (mn *Manager) ManageDataNoSQLtoSQL(postgr *SQLPostgres) error {
	err := mn.GetDataFromNoSQL()
	if err != nil {
		return err
	}
	return err
}

// GetSchemaFromNoSQL takes schema of tables from NoSQL entries
func (mn *Manager) GetSchemaFromNoSQL(keyspace string) error {
	mn.wg.Add(1)

	mn.SetKeyspace(keyspace)
	err := mn.cass.ReadTableSchema(keyspace, mn.posg.DbData)
	if err != nil {
		return err
	}

	defer mn.wg.Done()
	return err
}

// GetDataFromNoSQL reading data from NoSQL entries according to schema
func (mn *Manager) GetDataFromNoSQL() error {
	var err error
	for tabname, tab := range mn.cass.TableQueries {
		errchanel := make(chan error, 1)
		go func() {
			errchanel <- mn.ReaDataFromSingleNoSQLTable(tabname, tab.QueryInsert, tab.QuerySelect)
		}()
		err = <-errchanel
		if err != nil {
			return err
		}
	}
	mn.wg.Wait()
	return err
}

// ReaDataFromSingleNoSQLTable reads data from NoSQL table and write inti SQL table
func (mn *Manager) ReaDataFromSingleNoSQLTable(tabname string, insertquery string, selectquery string) error {
	var err error
	mn.wg.Add(1)
	tabcolums := mn.posg.DbData.Tables[tabname].Columns
	iter := mn.cass.CassandraSession.Query(selectquery).Iter()
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		errchanel := make(chan error, 1)
		mn.wg.Add(1)
		go func() {
			res := mn.ReturnSliceData(tabcolums, m)
			_, err = mn.posg.pdb.Exec(insertquery, res...)
			errchanel <- err
			mn.wg.Done()
		}()
		err = <-errchanel
		if err != nil {
			return err
		}
		m = map[string]interface{}{}
	}
	mn.wg.Done()
	return err
}

// ReturnSliceData return result from query ordering as sequence of fields in select query
func (mn *Manager) ReturnSliceData(tabcolums []Column, inMap map[string]interface{}) []interface{} {
	reschanel := make(chan []interface{}, 1)
	mn.wg.Add(1)
	values := []interface{}{}
	go func() {
		for i := 0; i < len(inMap); i++ {
			colname := tabcolums[i].Cname
			values = append(values, inMap[colname])

		}
		reschanel <- values
	}()
	mn.wg.Done()
	return <-reschanel
}
