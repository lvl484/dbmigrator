package main

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func Test_newSQLPostgres(t *testing.T) {
	tests := []struct {
		name    string
		want    *SQLPostgres
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newSQLPostgres()
			if (err != nil) != tt.wantErr {
				t.Errorf("newSQLPostgres() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSQLPostgres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLPostgres_DatabaseSQL(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	tests := []struct {
		name    string
		fields  fields
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			got, err := sp.DatabaseSQL()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.DatabaseSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLPostgres.DatabaseSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLPostgres_AddPrimaryKey(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		tabl string
		col  string
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			sp.AddPrimaryKey(tt.args.tabl, tt.args.col)
		})
	}
}

func TestSQLPostgres_AddForeignKey(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		constrname string
		forkey     ForeignKey
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			sp.AddForeignKey(tt.args.constrname, tt.args.forkey)
		})
	}
}

func TestSQLPostgres_AddColumnToTable(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		tname string
		col   Column
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			sp.AddColumnToTable(tt.args.tname, tt.args.col)
		})
	}
}

func TestSQLPostgres_ReadDBName(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			got, err := sp.ReadDBName()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.ReadDBName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLPostgres.ReadDBName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLPostgres_ReadTableSchema(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			if err := sp.ReadTableSchema(); (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.ReadTableSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLPostgres_ReadPrimaryKeys(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			if err := sp.ReadPrimaryKeys(); (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.ReadPrimaryKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLPostgres_ReadForeignKeys(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			if err := sp.ReadForeignKeys(); (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.ReadForeignKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLPostgres_ReadDataFromTable(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		querystring string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sql.Rows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			got, err := sp.ReadDataFromTable(tt.args.querystring)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.ReadDataFromTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLPostgres.ReadDataFromTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLPostgres_CreateNewDatabase(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		dbname string
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			if err := sp.CreateNewDatabase(tt.args.dbname); (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.CreateNewDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLPostgres_WriteSchemaToDatabase(t *testing.T) {
	type fields struct {
		pdb    *sql.DB
		DbData *DatabasePostg
	}
	type args struct {
		tabs map[string]TableQuery
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
			sp := &SQLPostgres{
				pdb:    tt.fields.pdb,
				DbData: tt.fields.DbData,
			}
			if err := sp.WriteSchemaToDatabase(tt.args.tabs); (err != nil) != tt.wantErr {
				t.Errorf("SQLPostgres.WriteSchemaToDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
