package goquery

import (
	"context"
)

type RdbmsDb interface {
	Connection() interface{}
	Transaction() (Tx, error)
	Select(dest any, tx *Tx, stmt string, params ...any) error
	SelectWithOptions(dest any, opts RdbmsQueryOptions) error
	Get(dest any, tx *Tx, stmt string, params ...any) error
	GetWithOptions(dest any, opts RdbmsQueryOptions) error
	Query(tx *Tx, stmt string, params ...any) (Rows, error)
	QueryWithOptions(opts RdbmsQueryOptions) (Rows, error)
	Insert(ds DataSet, rec any, tx *Tx) error
	InsertStmt(ds DataSet) (string, error)
	Exec(tx *Tx, stmt string, params ...any) error
	Execr(tx *Tx, stmt string, params ...any) (ExecResult, error)
	MustExec(tx *Tx, stmt string, params ...any)
	MustExecr(tx *Tx, stmt string, params ...any) ExecResult
	Batch() (Batch, error)
	SendBatch(batch Batch) BatchResult
	Close()
}

type RdbmsQueryOptions struct {
	Tx        *Tx
	Ctx       context.Context
	Statement string
	Params    []any
}
