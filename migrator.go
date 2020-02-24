package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// Running program without flags will migrate data from PostgreSQL to Cassandra
	// using a flag -cass will migrate data from Cassandra to PostgreSQL
	boolFromCassandra := flag.Bool("cass", false, "migrate from cassandra to postgresql")
	databaseCassandra := flag.String("db", "", "cassandra db to migrate")
	flag.Parse()

	if *boolFromCassandra {
		if strings.EqualFold(*databaseCassandra, "") {
			fmt.Printf("Please enter a name of database in command line to migrate from : example \" dbmigrator -cass -db mydatabasename\"")
			return
		}
		err := DoMigrateFromCassandra(*databaseCassandra)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		err := DoMigrateFromPostgres()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Thank you! Data were migrated PostgreSQL between Cassandra!")
}
