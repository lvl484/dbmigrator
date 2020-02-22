package main

import (
	"fmt"
)

func main() {
	var MigratorManager *Manager
	MigratorManager, err := NewManager()
	if err != nil {
		fmt.Printf("There is some problems with inicialization: %v\n", err)
		return
	}
	defer MigratorManager.CloseConnection()

	err = MigratorManager.GetSchemaFromSQL()
	if err != nil {
		fmt.Printf("Can not get shchema from PostgreSQL: %v\n", err)
		return
	}

	err = MigratorManager.PutSchemaToNoSQL()
	if err != nil {
		fmt.Printf("Can not put shchema to Cassandra: %v\n", err)
		return
	}
	err = MigratorManager.GetDataFromSQL()
	if err != nil {
		fmt.Printf("There is some problem with getting data from database: %v\n", err)
		return
	}
	err = MigratorManager.CloseConnection()
	if err != nil {
		fmt.Printf("Couldn't close connection to PostgreSQL database: %v\n", err)
		return
	}
	fmt.Println("Thank you! Data were migrated from PostgreSQL to Cassandra!")
}
