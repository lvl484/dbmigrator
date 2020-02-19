package main

var ConvTypePostgCasan = func() map[string]string {
	return map[string]string{
		"text":              "text",
		"character":         "text",
		"character varying": "varchar",
		"int":               "int",
		"boolean":           "boolean",
		"bool":              "boolean",
		"uuid":              "uuid",
		"inet":              "inet",
		"date":              "date",
		"double precision":  "double",
		"real":              "float",
		"timestamp":         "timestamp",
		"decimal":           "decimal",
		"serial":            "counter",
	}
}

type SQLQueries string

const (
	DbNameQuery SQLQueries = "SELECT * FROM information_schema.information_schema_catalog_name"
	TablesQuery SQLQueries = "SELECT tb.table_name, ordinal_position, column_name, data_type FROM information_schema.tables tb " +
		"JOIN information_schema.columns cs ON tb.table_name = cs.table_name WHERE tb.table_type = 'BASE TABLE' " +
		"AND tb.table_schema = 'public' AND tb.table_catalog = $1 order by tb.table_name, ordinal_position;"
	PrimaryKeysQuery SQLQueries = "SELECT tc.table_name,column_name FROM information_schema.table_constraints  tc " +
		"join information_schema.constraint_column_usage cu on tc.constraint_name = cu.constraint_name " +
		"where tc.table_schema = 'public' and constraint_type = 'PRIMARY KEY'and tc.table_catalog = $1;"
	ForeignKeysQuery SQLQueries = "select tc.constraint_name, tc.table_name, kcu.column_name, ccu.table_name AS foreign_table_name, " +
		"ccu.column_name AS foreign_column_name from information_schema.table_constraints AS tc " +
		"JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name " +
		"JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name " +
		"where tc.table_schema = 'public' and constraint_type = 'FOREIGN KEY' and tc.table_catalog = $1;"
	DataTablesQuery SQLQueries = "SELECT * FROM $1"
)

type NoSQLQueries string

const (
	CreateKeyspaceQuery NoSQLQueries = "CREATE KEYSPACE %s WITH replication = { class : \"SimpleStrategy\", \"replication_factor\" : %d }"
	DropKeyspaceQuery   NoSQLQueries = "DROP KEYSPACE IF EXISTS %s"

//	MakeForeignKeysQuery   NoSQLQueries = "CREATE TABLE Dependencies (ForeignKey text PRIMARY KEY, MainTable text, MainColumn text, DependTable text, DependColumn)"
//	InsertForeignKeysQuery NoSQLQueries = "COPY Dependencies (ForeignKey text, MainTable text, MainColumn text, DependTable text, DependColumn) FROM STDIN"
)
