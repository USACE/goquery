package oracle

import (
	"fmt"

	"github.com/usace/goquery/v3"
)

const (
	registryName string = "godror"
)

func init() {
	goquery.DbRegistry[registryName] = OracleDialect
}

var OracleDialect = goquery.DbDialect{
	TableExistsStmt: `select count(*) from user_tables where table_name=$1`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf(":%s", field)
	},
	Seq: func(sequence string) string {
		return fmt.Sprintf("nextval('%s')", sequence)
	},
	Url: func(config *goquery.RdbmsConfig) string {
		if config.OnInit == "" {
			return fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s" libDir="%s"`,
				config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname, config.ExternalLib)
		}
		return fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s" libDir="%s" onInit="%s" %s`,
			config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname, config.ExternalLib, config.OnInit, config.DbDriverSettings)
	},
}
