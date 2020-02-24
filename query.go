package main

var ConvTypePostgCasan = func() map[string]string {
	return map[string]string{
		`text`:                        `text`,
		`character`:                   `text`,
		`character varying`:           `varchar`,
		`int`:                         `int`,
		`integer`:                     `int`,
		`boolean`:                     `boolean`,
		`bool`:                        `boolean`,
		`uuid`:                        `uuid`,
		`inet`:                        `inet`,
		`date`:                        `date`,
		`double precision`:            `double`,
		`real`:                        `double`,
		`timestamp`:                   `timestamp`,
		`decimal`:                     `decimal`,
		`serial`:                      `counter`,
		`timestamp without time zone`: `timestamp`,
		`smallint`:                    `smallint`,
		`numeric`:                     `double`,
		`USER-DEFINED`:                `text`,
		`ARRAY`:                       `text`,
		`tsvector`:                    `text`,
		`bytea`:                       `blob`,
	}
}

var ConvTypeCasanPost = func() map[string]string {
	return map[string]string{
		`ascii`: `text`, `bigint`: `bigint`, `blob`: `bytea`, `boolean`: `boolean`, `counter`: `int`, `date`: `date`, `decimal`: `real`, `double`: `real`, `duration`: `text`, `float`: `real`, `inet`: `inet`, `int`: `int`, `smallint`: `smallint`, `text`: `text`, `time`: `time`, `timestamp`: `timestamp`, `timeuuid`: `uuid`, `tinyint`: `smallint`, `uuid`: `uuid`, `varchar`: `text`, `varint`: `int`,
	}
}

const (
	PostgresConfDB    = `host=%s port=%s user=%s password=%s dbname=%s`
	DbNameQuery       = `SELECT * FROM information_schema.information_schema_catalog_name`
	TablesSchemaQuery = `SELECT tb.table_name, ordinal_position, column_name, data_type FROM information_schema.tables tb 
		JOIN information_schema.columns cs ON tb.table_name = cs.table_name WHERE tb.table_type = 'BASE TABLE'
		AND tb.table_schema = 'public' AND tb.table_catalog = $1 order by tb.table_name, ordinal_position;`
	PrimaryKeysQuery = `SELECT tc.table_name,column_name FROM information_schema.table_constraints  tc 
		join information_schema.constraint_column_usage cu on tc.constraint_name = cu.constraint_name 
		where tc.table_schema = 'public' and constraint_type = 'PRIMARY KEY' and tc.table_catalog = $1;`
	ForeignKeysQuery = `SELECT tc.constraint_name, tc.table_name, kcu.column_name, ccu.table_name AS foreign_table_name, 
		ccu.column_name AS foreign_column_name from information_schema.table_constraints AS tc 
		JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name 
		JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name 
		where tc.table_schema = 'public' and constraint_type = 'FOREIGN KEY' and tc.table_catalog = $1;`
	PostgresSelect        = `SELECT %s FROM %s`
	KeyspaceQuery         = `CREATE KEYSPACE IF NOT EXISTS %s WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`
	CassandraTable        = `CREATE TABLE IF NOT EXISTS %s.%s (%s)`
	PostgreSQLTable       = `CREATE TABLE %s (%s)`
	CassandraPrimary      = ` PRIMARY KEY (%s)`
	CassandraCopyData     = `INSERT INTO %s.%s (%s) VALUES (%s) IF NOT EXISTS`
	PostgreSQLInsertData  = `INSERT INTO %s (%s) VALUES (%s)`
	CassandraTableSchema  = `SELECT table_name FROM system_schema.tables WHERE keyspace_name = '%s';`
	CassandraColumnSchema = `SELECT column_name, type, kind FROM system_schema.columns 
		WHERE keyspace_name = '%s' AND table_name = '%s';`
)

// CompareType function that check postgresql data type that would be converted
// to type double, type without precision to migrate in Cassandra DB
func CompareType(coltype string) bool {
	ok := false
	switch coltype {
	case `double precision`:
		ok = true
	case `numeric`:
		ok = true
	case `real`:
		ok = true
	case `decimal`:
		ok = true
	}

	return ok
}
