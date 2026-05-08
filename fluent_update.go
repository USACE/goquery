package goquery

import (
	"fmt"
	"log"
)

// FluentUpdate provides a fluent interface for UPDATE operations
// It follows the same pattern as FluentSelect, using DataSet and StatementKey
type FluentUpdate struct {
	store DataStore
	ds    DataSet
	tx    *Tx
	ui    UpdateInput
}

// UpdateInput contains the parameters for an UPDATE operation
type UpdateInput struct {
	StatementKey string
	Statement    string
	Suffix       string
	BindParams   []interface{}
	StmtAppends  []interface{}
	PanicOnErr   bool
	LogSql       bool
}

// DataSet sets the DataSet containing the UPDATE statement
func (u *FluentUpdate) DataSet(ds DataSet) *FluentUpdate {
	u.ds = ds
	return u
}

// StatementKey specifies which statement to use from the DataSet's Statements map
func (u *FluentUpdate) StatementKey(key string) *FluentUpdate {
	u.ui.StatementKey = key
	return u
}

// Apply adds values to be appended to the statement using fmt.Sprintf
// WARNING: Only use with internal values, never user input (SQL injection risk)
func (u *FluentUpdate) Apply(vals ...interface{}) *FluentUpdate {
	u.ui.StmtAppends = vals
	return u
}

// Params sets the bind parameters for the UPDATE statement
func (u *FluentUpdate) Params(params ...interface{}) *FluentUpdate {
	u.ui.BindParams = params
	return u
}

// Suffix adds additional SQL to append to the statement (e.g., RETURNING clause)
func (u *FluentUpdate) Suffix(suffix string) *FluentUpdate {
	u.ui.Suffix = suffix
	return u
}

// Tx sets the transaction to use for this update
func (u *FluentUpdate) Tx(tx *Tx) *FluentUpdate {
	u.tx = tx
	return u
}

// PanicOnErr sets whether to panic on error (useful in transactions)
func (u *FluentUpdate) PanicOnErr(panicOnErr bool) *FluentUpdate {
	u.ui.PanicOnErr = panicOnErr
	return u
}

// LogSql enables SQL logging for debugging
func (u *FluentUpdate) LogSql(logSql bool) *FluentUpdate {
	u.ui.LogSql = logSql
	return u
}

// Execute runs the UPDATE statement and returns an error if one occurs
func (u *FluentUpdate) Execute() error {
	stmt, err := u.getStatement()
	if err != nil {
		if u.ui.PanicOnErr {
			panic(err)
		}
		return err
	}

	if u.ui.LogSql {
		log.Println(stmt, u.ui.BindParams)
	}

	err = u.store.Exec(u.tx, stmt, u.ui.BindParams...)
	if err != nil && u.ui.PanicOnErr {
		panic(err)
	}

	return err
}

// Execr executes the UPDATE and returns the result with rows affected
func (u *FluentUpdate) Execr() (ExecResult, error) {
	stmt, err := u.getStatement()
	if err != nil {
		if u.ui.PanicOnErr {
			panic(err)
		}
		return nil, err
	}

	if u.ui.LogSql {
		log.Println(stmt, u.ui.BindParams)
	}

	result, err := u.store.Execr(u.tx, stmt, u.ui.BindParams...)
	if err != nil && u.ui.PanicOnErr {
		panic(err)
	}

	return result, err
}

// getStatement retrieves the UPDATE statement from the DataSet or uses the provided statement
// Follows the same pattern as getSelectStatement in sql_generators.go
func (u *FluentUpdate) getStatement() (string, error) {
	var stmt string
	var ok bool

	// Get statement from DataSet using key, or use provided statement
	if u.ui.StatementKey != "" || u.ui.Statement == "" {
		if u.ds == nil {
			return "", fmt.Errorf("missing dataset when referencing a statement key")
		}
		if u.ui.StatementKey != "" {
			if stmt, ok = u.ds.Commands()[u.ui.StatementKey]; !ok {
				return "", fmt.Errorf("unable to find statement for %s: %s", u.ds.Entity(), u.ui.StatementKey)
			}
		}
	} else {
		stmt = u.ui.Statement
	}

	// Apply suffix and statement appends using fmt.Sprintf (same as FluentSelect)
	stmt = fmt.Sprintf("%s %s", stmt, u.ui.Suffix)
	return fmt.Sprintf(stmt, u.ui.StmtAppends...), nil
}

type updateError struct {
	msg string
}

func (e *updateError) Error() string {
	return e.msg
}
