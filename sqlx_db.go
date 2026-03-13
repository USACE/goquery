package goquery

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/jmoiron/sqlx"
)

type SqlRows struct {
	rows       *sql.Rows
	rowScanner *sqlscan.RowScanner
}

type SqlxExecResult struct {
	res sql.Result
}

func (ser SqlxExecResult) RowsAffected() int64 {
	rows, err := ser.res.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	return rows
}

func (s *SqlRows) Columns() ([]string, error) {
	return s.rows.Columns()
}

func (s *SqlRows) ColumnTypes() ([]reflect.Type, error) {
	sts, err := s.rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	t := make([]reflect.Type, len(sts))
	for i := range sts {
		t[i] = sts[i].ScanType()
	}
	return t, nil
}

func (s *SqlRows) Next() bool {
	return s.rows.Next()
}

func (s *SqlRows) Scan(dest ...interface{}) error {
	return s.rows.Scan(dest...)
}

func (s *SqlRows) ScanStruct(dest interface{}) error {
	if s.rowScanner == nil {
		s.rowScanner = sqlscan.NewRowScanner(s.rows)
	}
	return s.rowScanner.Scan(dest)
}

func (s *SqlRows) Close() error {
	return s.rows.Close()
}

type SqlxDb struct {
	db      *sqlx.DB
	dialect DbDialect
}

func getDialect(driver string) (DbDialect, error) {
	switch driver {
	case "pgx":
		return pgDialect, nil
	case "godror":
		return oracleDialect, nil
	case "sqlite":
		return sqliteDialect, nil
	default:
		return DbDialect{}, errors.New(fmt.Sprintf("Unsupported DB Driver: %s", driver))
	}
}

func NewSqlxConnection(config *RdbmsConfig) (SqlxDb, error) {
	dialect, err := getDialect(config.DbDriver)
	if err != nil {
		return SqlxDb{}, err
	}
	dburl := dialect.Url(config)
	con, err := sqlx.Connect(config.DbDriver, dburl)

	return SqlxDb{con, dialect}, err
}

func (sdb *SqlxDb) querier(tx *Tx) sqlx.Queryer {
	if tx != nil {
		return tx.SqlXTx()
	}
	return sdb.db
}

type sqlxexecr interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (sdb *SqlxDb) execr(tx *Tx) sqlxexecr {
	if tx != nil {
		return tx.SqlTx()
	}
	return sdb.db
}

func (sdb *SqlxDb) Connection() interface{} {
	return sdb.db
}

func (sdb *SqlxDb) Select(dest interface{}, tx *Tx, stmt string, params ...interface{}) error {
	if len(params) == 0 {
		return sqlx.Select(sdb.querier(tx), dest, stmt)
	}
	return sqlx.Select(sdb.querier(tx), dest, stmt, params...)
}

func (sdb *SqlxDb) Get(dest interface{}, tx *Tx, stmt string, params ...interface{}) error {
	if len(params) == 0 {
		return sqlx.Get(sdb.querier(tx), dest, stmt)
	}
	return sqlx.Get(sdb.querier(tx), dest, stmt, params...)
}

func (sdb *SqlxDb) Query(tx *Tx, stmt string, params ...interface{}) (Rows, error) {
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.SqlTx().Query(stmt, params...)
	} else {
		rows, err = sdb.db.Query(stmt, params...)
	}
	return &SqlRows{rows, nil}, err
}

func (sdb *SqlxDb) Exec(tx *Tx, stmt string, params ...interface{}) error {
	_, err := sdb.execr(tx).Exec(stmt, params...)
	return err
}

func (sdb *SqlxDb) Execr(tx *Tx, stmt string, params ...interface{}) (ExecResult, error) {
	res, err := sdb.execr(tx).Exec(stmt, params...)
	return SqlxExecResult{res}, err
}

func (sdb *SqlxDb) MustExec(tx *Tx, stmt string, params ...interface{}) {
	if tx != nil {
		tx.SqlXTx().MustExec(stmt, params...)
	} else {
		sdb.db.MustExec(stmt, params...)
	}
}

func (sdb *SqlxDb) MustExecr(tx *Tx, stmt string, params ...interface{}) ExecResult {
	var res sql.Result
	if tx != nil {
		res = tx.SqlXTx().MustExec(stmt, params...)
	} else {
		res = sdb.db.MustExec(stmt, params...)
	}
	return SqlxExecResult{res}
}

func (sdb *SqlxDb) Batch() (Batch, error) {
	return nil, errors.New("batch operations are not supported by the sqlx driver")
}

func (sdb *SqlxDb) SendBatch(batch Batch) BatchResult {
	return nil
}

func (sdb *SqlxDb) InsertStmt(ds DataSet) (string, error) {
	return ToInsert(ds, sdb.dialect)
}

func (sdb *SqlxDb) Insert(ds DataSet, rec interface{}, tx *Tx) error {
	var err error
	var stmt string
	var ok bool

	if stmt, ok = ds.Commands()["insert"]; !ok {
		stmt, err = ToInsert(ds, sdb.dialect)
	}
	if err != nil {
		return err
	}
	params := StructToIArray(rec)
	if tx == nil {
		_, err = sdb.db.Exec(stmt, params...)
		return err
	} else {
		sqltx := tx.SqlTx()
		_, err = sqltx.Exec(stmt, params...)
		return err
	}
}

func (sdb *SqlxDb) Transaction() (Tx, error) {
	tx, err := sdb.db.Beginx()
	return Tx{tx}, err
}
