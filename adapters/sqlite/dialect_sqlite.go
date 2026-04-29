package sqlite

import "github.com/usace/goquery/v3"

const (
	registryNameCgo      string = "sqlite3"
	registryNameNativeGo string = "sqlite"
)

func init() {
	goquery.DbRegistry[registryNameCgo] = SqliteDialect
	goquery.DbRegistry[registryNameNativeGo] = SqliteDialect
}

var SqliteDialect = goquery.DbDialect{
	TableExistsStmt: `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`,
	Bind: func(field string, i int) string {
		return "?"
	},
	Seq: func(sequence string) string {
		//sequences are not supported in sqlite
		return ""
	},
	Url: func(config *goquery.RdbmsConfig) string {
		return config.Dbname
	},
}
