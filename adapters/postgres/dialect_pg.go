package postgres

import (
	"fmt"

	"github.com/usace/goquery/v3"
)

const (
	registryName string = "pgx"
)

func init() {
	goquery.DbRegistry[registryName] = PgDialect
}

var PgDialect = goquery.DbDialect{
	TableExistsStmt: `SELECT count(*) FROM information_schema.tables WHERE  table_schema = $1 AND table_name = $2`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf("$%d", i+1)
	},
	Seq: func(sequence string) string {
		return fmt.Sprintf("nextval('%s')", sequence)
	},

	//only used by sqlx.  pgx url is constructed directly in the NewPgxConnection function in pgx_dg.go
	Url: func(config *goquery.RdbmsConfig) string {
		return fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s",
			config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname)
	},
}
