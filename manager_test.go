package main

import (
	"reflect"
	"sync"
	"testing"
)

type DataAccessLayer interface {

  }

type TestData {

}

func TestNewManager(t *testing.T) {
	type args struct {
		ver int
	}
	tests := []struct {
		name    string
		args    args
		want    *Manager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager(tt.args.ver)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetKeyspace(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		keyspasename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			mn.SetKeyspace(tt.args.keyspasename)
		})
	}
}

func TestManager_GetSchemaFromSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.GetSchemaFromSQL(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSchemaFromSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_PutSchemaToNoSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.PutSchemaToNoSQL(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.PutSchemaToNoSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetDataFromSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.GetDataFromSQL(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetDataFromSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ReaDataFromSingleTable(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		selectquery string
		insertquery string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.ReaDataFromSingleTable(tt.args.selectquery, tt.args.insertquery); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ReaDataFromSingleTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_CloseConnection(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			mn.CloseConnection()
		})
	}
}

func TestDoMigrateFromPostgres(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoMigrateFromPostgres(); (err != nil) != tt.wantErr {
				t.Errorf("DoMigrateFromPostgres() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoMigrateFromCassandra(t *testing.T) {
	type args struct {
		cassName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoMigrateFromCassandra(tt.args.cassName); (err != nil) != tt.wantErr {
				t.Errorf("DoMigrateFromCassandra() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ManageDataNoSQLtoSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		postgr *SQLPostgres
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.ManageDataNoSQLtoSQL(tt.args.postgr); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ManageDataNoSQLtoSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetSchemaFromNoSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		keyspace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.GetSchemaFromNoSQL(tt.args.keyspace); (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSchemaFromNoSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetDataFromNoSQL(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.GetDataFromNoSQL(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetDataFromNoSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ReaDataFromSingleNoSQLTable(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		tabname     string
		insertquery string
		selectquery string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if err := mn.ReaDataFromSingleNoSQLTable(tt.args.tabname, tt.args.insertquery, tt.args.selectquery); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ReaDataFromSingleNoSQLTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ReturnSliceData(t *testing.T) {
	type fields struct {
		posg *SQLPostgres
		cass *Cassandra
		wg   sync.WaitGroup
	}
	type args struct {
		tabcolums []Column
		inMap     map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn := &Manager{
				posg: tt.fields.posg,
				cass: tt.fields.cass,
				wg:   tt.fields.wg,
			}
			if got := mn.ReturnSliceData(tt.args.tabcolums, tt.args.inMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.ReturnSliceData() = %v, want %v", got, tt.want)
			}
		})
	}
}
