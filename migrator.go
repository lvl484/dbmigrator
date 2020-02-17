package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome in DB SQLtoNOSQL migrator!")
	inputPostg := NewSQLPostgre()
	inputPostg.GetSchemaFromSQl()
	inputPostg.PutSchemaToNoSQL()
	inputPostg.GetDataFromSQL()
}
