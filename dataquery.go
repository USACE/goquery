package goquery

import (
	"errors"
	"fmt"
	"reflect"
)

const selectkey = "select"

//const updatekey = "update"
//const insertkey = "insert"

type RowFunction func(r Rows) error

type Rows interface {
	Columns() ([]string, error)
	ColumnTypes() ([]reflect.Type, error)
	Next() bool
	Scan(dest ...interface{}) error
	ScanStruct(dest interface{}) error
	Close() error
}

type DataSet interface {
	Entity() string
	FieldSlice() interface{} //@depricated.  Will be removed in the next version...maybe
	Fields() interface{}     //@depricated.  Will be removed in the next version...maybe
	Commands() map[string]string
	PutCommand(key string, stmt string)
}

type Statements map[string]string

func (s Statements) Get(key string) (string, error) {
	if val, ok := s[key]; ok {
		return val, nil
	}
	return "", errors.New("Invalid statement")
}
func (s Statements) GetOrPanic(key string) string {
	if val, ok := s[key]; ok {
		return val
	}
	panic(errors.New("Invalid statement"))
}

type TableDataSet struct {
	Name        string
	Schema      string //optional
	Statements  Statements
	TableFields any
}

func (t *TableDataSet) FieldSlice() interface{} {
	typ := reflect.TypeOf(t.TableFields)
	slice := reflect.New(reflect.SliceOf(typ))
	return slice.Interface()
}

func (t *TableDataSet) Fields() interface{} {
	return t.TableFields
}

func (t *TableDataSet) Entity() string {
	if t.Schema != "" {
		return fmt.Sprintf("%s.%s", t.Schema, t.Name)
	}
	return t.Name
}

func (t *TableDataSet) Commands() map[string]string {
	return t.Statements
}

func (t *TableDataSet) PutCommand(key string, stmt string) {
	if t.Statements == nil {
		t.Statements = make(map[string]string)
	}
	t.Statements[selectkey] = stmt
}
