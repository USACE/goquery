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
	ToMap() (map[string]any, error)
	Close() error
}

// converts the current Rows position to a map
func RowToMap(r Rows) (map[string]any, error) {
	cols, err := r.Columns()
	if err != nil {
		return nil, err
	}

	colTypes, err := r.ColumnTypes()
	if err != nil {
		return nil, err
	}

	vals := make([]any, len(cols))

	for i := range vals {
		pval := reflect.New(colTypes[i])
		ival := pval.Interface() //call Elem to dereference the pointer created by reflect.New
		vals[i] = ival
	}

	err = r.Scan(vals...)
	valmap := make(map[string]any)
	//this is pretty gross, but it is significantly faster than reflection which is the fallback
	for i, col := range cols {
		val := vals[i]
		var concreteVal any

		switch v := val.(type) {
		case *string:
			concreteVal = *v
		case *int64:
			concreteVal = *v
		case *int32:
			concreteVal = *v
		case *float64:
			concreteVal = *v
		case *float32:
			concreteVal = *v
		case *bool:
			concreteVal = *v
		default:
			// Fallback to reflection ONLY for unknown types
			concreteVal = reflect.Indirect(reflect.ValueOf(val)).Interface()
		}

		valmap[col] = concreteVal
	}

	return valmap, nil
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
	t.Statements[key] = stmt
}
