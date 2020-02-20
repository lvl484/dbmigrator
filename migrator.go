package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

const cCBootstrapTimeout time.Duration = 3 * time.Second

var CassandraSession *gocql.Session
var DBKeyspace string
var CassandraCluster *gocql.ClusterConfig

func main() {
	var MigratorManager *Manager

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
	CassandraCluster := gocql.NewCluster(HOST)
	CassandraCluster.Timeout = cCBootstrapTimeout
	CassandraSession, err = CassandraCluster.CreateSession()
	if err != nil {
		fmt.Println("Can not connect Cassandra")
		return
	}
	defer CassandraSession.Close()
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
