package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

var CassandraSession *gocql.Session
var DBKeyspace string
var CassandraCluster *gocql.ClusterConfig

func main() {
	var MigratorManager *Manager
	defer CassandraSession.Close()
	MigratorManager, err := NewManager()
	if err != nil {
		fmt.Println("There is some problem with connection to database")
		return
	}
	err = MigratorManager.GetSchemaFromSQl()
	if err != nil {
		fmt.Println("Can not get shchema from PostgreSQL")
		return
	}

	err = MigratorManager.PutSchemaToNoSQL()
	if err != nil {
		fmt.Println("Can not put shchema to Cassandra")
		return
	}
	err = MigratorManager.GetDataFromSQL()
	if err != nil {
		fmt.Println("There is some problem with getting data from database")
		return
	}
	MigratorManager.wg.Wait()

}
