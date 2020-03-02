package main

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"
)

func TestNewCassandra(t *testing.T) {
	tests := []struct {
		name    string
		want    *Cassandra
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCassandra()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCassandra() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCassandra() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCassandra_CreateKeyspaceCassandra(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			if err := cs.CreateKeyspaceCassandra(); (err != nil) != tt.wantErr {
				t.Errorf("Cassandra.CreateKeyspaceCassandra() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCassandra_CreateNoSQLTableScheme(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		dbD *DatabasePostg
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			cs.CreateNoSQLTableScheme(tt.args.dbD)
		})
	}
}

func TestCassandra_CreateQueryForSQLTable(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		tablename string
		tab       Table
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			cs.CreateQueryForSQLTable(tt.args.tablename, tt.args.tab)
		})
	}
}

func TestCassandra_WriteSchemaToDB(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			if err := cs.WriteSchemaToDB(); (err != nil) != tt.wantErr {
				t.Errorf("Cassandra.WriteSchemaToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCassandra_ReadTableSchema(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		keyspname string
		datapos   *DatabasePostg
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			if err := cs.ReadTableSchema(tt.args.keyspname, tt.args.datapos); (err != nil) != tt.wantErr {
				t.Errorf("Cassandra.ReadTableSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCassandra_ReadColumnSchema(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		keyspname string
		tablename string
		datapos   *DatabasePostg
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			if err := cs.ReadColumnSchema(tt.args.keyspname, tt.args.tablename, tt.args.datapos); (err != nil) != tt.wantErr {
				t.Errorf("Cassandra.ReadColumnSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCassandra_CreateSQLTableScheme(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		dbD *DatabasePostg
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			if err := cs.CreateSQLTableScheme(tt.args.dbD); (err != nil) != tt.wantErr {
				t.Errorf("Cassandra.CreateSQLTableScheme() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCassandra_CreateQueryForNoSQLTable(t *testing.T) {
	type fields struct {
		CassandraSession *gocql.Session
		DBKeyspace       string
		CassandraCluster *gocql.ClusterConfig
		TableQueries     map[string]TableQuery
	}
	type args struct {
		tablename string
		tab       Table
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
			cs := &Cassandra{
				CassandraSession: tt.fields.CassandraSession,
				DBKeyspace:       tt.fields.DBKeyspace,
				CassandraCluster: tt.fields.CassandraCluster,
				TableQueries:     tt.fields.TableQueries,
			}
			cs.CreateQueryForNoSQLTable(tt.args.tablename, tt.args.tab)
		})
	}
}
