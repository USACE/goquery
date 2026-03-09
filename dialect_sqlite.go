package goquery

var sqliteDialect = DbDialect{
	TableExistsStmt: `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`,
	Bind: func(field string, i int) string {
		return "?"
	},
	Seq: func(sequence string) string {
		//sequences are not supported in sqlite
		return ""
	},
	Url: func(config *RdbmsConfig) string {
		return config.Dbname
	},
}
