package duckdb

import (
	"fmt"

	"github.com/usace/goquery/v3"
)

const (
	registryName string = "duckdb"
)

func init() {
	goquery.DbRegistry[registryName] = DuckdbDialect
}

var DuckdbDialect = goquery.DbDialect{
	TableExistsStmt: `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf("$%d", i+1)
	},
	Url: func(config *goquery.RdbmsConfig) string {
		return config.Dbname
	},
}
