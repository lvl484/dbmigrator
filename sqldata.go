package main

type Column struct {
	Cname string
	Ctype string
	Pk    bool
}

type Table struct {
	Columns []Column
}

type ForeignKey struct {
	Tablename string
	Colname   string
	Tableref  string
	Colref    string
}

type DatabasePostg struct {
	Databasename string
	Tables       map[string]Table
	Foreignkeys  map[string]ForeignKey
}

// MakePrimaryKey change column to primary key
func (c *Column) MakePrimaryKey() {
	c.Pk = true
}

// SetDBName set a name of DB to our structure
func (ds *DatabasePostg) SetDBName(name string) {
	ds.Databasename = name
}

// NewDatabasePostg create a new instance of DatabasePostg
func NewDatabasePostg() *DatabasePostg {
	return &DatabasePostg{
		Databasename: "",
		Tables:       make(map[string]Table),
		Foreignkeys:  make(map[string]ForeignKey),
	}
}
