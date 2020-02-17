package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

const HOST = "localhost"

type Cassandra struct {
	clName      string
	TableInsert map[string]string
	clust       *gocql.ClusterConfig
	csCQL       map[int]NoSQLQueries
}

func CreateCassandra(dbname string) *Cassandra {
	cluster := gocql.NewCluster(HOST)
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	err := session.Query(`DROP KEYSPACE IF EXISTS ` + dbname).Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = session.Query(fmt.Sprintf(`CREATE KEYSPACE %s
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : %d
    }`, dbname, 1)).Exec()

	if err != nil {
		log.Fatal(err)
	}

	csq := initCassQuery()
	cluster.Keyspace = dbname
	mm := make(map[string]string)

	return &Cassandra{
		clName:      dbname,
		TableInsert: mm,
		clust:       cluster,
		csCQL:       csq,
	}
}

func initCassQuery() map[int]NoSQLQueries {
	q := make(map[int]NoSQLQueries, 4)
	q[0] = CreateTables
	q[1] = InsertDataTable
	q[2] = MakeForeignKeys
	q[3] = InsertForeignKeys
	return q
}

func (cs *Cassandra) InputData(tableName string, chIn <-chan *sql.Rows) {
	session, _ := cs.clust.CreateSession()
	defer session.Close()
	var srows *sql.Rows
	srows = <-chIn
	for srows.Next() {
		if err := session.Query(cs.TableInsert[tableName], srows.Scan()).Exec(); err != nil {
			log.Println(err)
		}
	}

}

func (cs *Cassandra) RunQuery(order int, parnum int, chIn <-chan []string) {
	session, _ := cs.clust.CreateSession()
	defer session.Close()

	var ss, ds []string
	ss = <-chIn
	qs, is := cs.PrepareQS(order, parnum, len(ss))
	ds = cs.PrepareData(order, parnum, ss, is)
	if err := session.Query(qs, ds).Exec(); err != nil {
		log.Println(err)
	}
}

func (cs *Cassandra) PrepareData(order int, parnum int, ss []string, inpstr string) []string {
	var newss []string
	newss = append(newss, ss[0])
	instr := fmt.Sprintf("INSERT INTO %s (", ss[0])
	for i := 1; i < len(ss)-parnum; i++ {
		newss = append(newss, fmt.Sprintf("%s %s", ss[i], ConvTypePostgCasan()[ss[i+1]]))
		if i == (len(ss) - parnum - 2) {
			instr += fmt.Sprintf(" %s", ss[i])
		} else {
			instr += fmt.Sprintf(" %s,", ss[i])
		}
		i++
	}
	instr += fmt.Sprintf(") VALUES (%s) FROM STDIN", inpstr)
	cs.TableInsert[ss[0]] = instr
	for i := len(ss) - parnum; i < len(ss); i++ {
		newss = append(newss, ss[i])
	}
	return newss
}

func (cs *Cassandra) PrepareQS(order int, parnum int, countpar int) (string, string) {
	ss := ""
	var ins, si string
	switch order {
	case 1:
		i := countpar - 1 - parnum
		si, ins = retSQLstring(i/2, parnum)
		ss = string(cs.csCQL[0]) + si
	case 2:
	}
	return ss, ins
}

func retSQLstring(ct int, pr int) (string, string) {
	ss := ""
	is := ""
	for i := 0; i < ct; i++ {
		if i == (ct - 1) {
			ss += "? ?"
			is += " ?"
		} else {
			ss += "? ?,"
			is += " ?, "
		}
	}
	for i := 0; i < pr; i++ {
		ss += ", ? "
	}
	ss += ")"
	return ss, is
}
