package goquery

import "fmt"

var duckdbDialect = DbDialect{
	TableExistsStmt: `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf("$%d", i+1)
	},
}
