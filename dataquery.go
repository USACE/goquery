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
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
	ScanStruct(dest interface{}) error
	ToMap() (map[string]any, error)
	ScanIntoMap(cols []string, targetMap map[string]any) error
	Close() error
}

// RowBuffers holds the pre-allocated slices needed for a dynamic scan.
type RowBuffers struct {
	Vals     []any
	ScanDest []any
}

// RowToMapPool parses a dynamic row with near-zero memory allocations.
// To maximize performance, pass an existing map to reuse (targetMap).
func RowToMapBuf(r Rows, cols []string, targetMap map[string]any, buf *RowBuffers) error {
	length := len(cols)

	// Map pointers to the target slice indices
	for i := 0; i < length; i++ {
		buf.Vals[i] = nil // Reset old data
		buf.ScanDest[i] = &buf.Vals[i]
	}

	// 3. Scan directly into the pooled pointers
	if err := r.Scan(buf.ScanDest...); err != nil {
		return err
	}

	// 4. Map the columns into your target map without reflection
	for i, col := range cols {
		rawVal := buf.Vals[i]

		if rawVal == nil {
			targetMap[col] = nil
			continue
		}

		switch v := rawVal.(type) {
		case string:
			targetMap[col] = v
		case int64:
			targetMap[col] = v
		case float64:
			targetMap[col] = v
		case bool:
			targetMap[col] = v
		case []byte:
			// Treat byte arrays safely. If it's a string/geometry representation, cast it.
			targetMap[col] = string(v)
		default:
			targetMap[col] = v
		}
	}

	return nil
}

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
		// FIX: Create a pointer to the type first.
		// If colTypes[i] is 'string', this makes it '*string'.
		ptrType := reflect.PointerTo(colTypes[i])

		// reflect.New(ptrType) creates a pointer to that type.
		// So we get '**string'. This double-pointer allows 'pgx' to set it to nil on NULL.
		pval := reflect.New(ptrType)
		vals[i] = pval.Interface()
	}

	// Now scanning into **Type (e.g., **string)
	err = r.Scan(vals...)
	if err != nil {
		return nil, err
	}

	valmap := make(map[string]any)

	for i, col := range cols {
		// val is **string (or **int, etc.)
		val := vals[i]

		// Dereference the outer pointer to get the inner pointer (*string)
		// We use reflection here because we don't know the specific type yet
		outerPtr := reflect.ValueOf(val)
		innerPtr := outerPtr.Elem() // This is the *string

		// CHECK FOR NULL: If the inner pointer is nil, the DB value was NULL
		if innerPtr.IsNil() {
			valmap[col] = nil
			continue
		}

		// If not nil, we extract the actual *string (or *int) to pass to your switch
		actualPtr := innerPtr.Interface()

		var concreteVal any

		switch v := actualPtr.(type) {
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
			// Fallback: dereference the pointer we verified is not nil
			concreteVal = reflect.Indirect(reflect.ValueOf(actualPtr)).Interface()
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
