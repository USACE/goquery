package goquery

import (
	"reflect"
	//"github.com/ulule/deepcopier"
)

func TagsAndVals(tag string, data interface{}) ([]string, []interface{}) {
	val := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()
	fieldNum := val.NumField()
	tags := make([]string, fieldNum)
	ia := make([]interface{}, fieldNum)
	//@TODO add recursion for type encapsulation
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup(tag); ok {
			tags[i] = tagval
		}
		ia[i] = val.Field(i).Addr().Interface()
	}
	return tags, ia
}

func ValueMap(tag string, data interface{}) map[string]interface{} {
	val := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()
	fieldNum := val.NumField()
	valmap := make(map[string]interface{})
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup(tag); ok {
			valmap[tagval] = val.Field(i).Addr().Interface()

		}
	}
	return valmap
}

func TagAsPositionMap(tag string, data interface{}) map[string]int {
	tagmap := make(map[string]int)
	typ := reflect.TypeOf(data).Elem()
	fieldNum := typ.NumField()
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup(tag); ok {
			tagmap[tagval] = i
		}
	}
	return tagmap
}

func TagAsStringArray(tag string, data interface{}) []string {
	typ := reflect.TypeOf(data).Elem()
	fieldNum := typ.NumField()
	tags := make([]string, fieldNum)
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup(tag); ok {
			tags[i] = tagval
		}
	}
	return tags
}

func isSlice(data interface{}) bool {
	rval := reflect.ValueOf(data)
	val := reflect.Indirect(rval)
	return val.Kind() == reflect.Slice
}

func StructToIArray(data interface{}) []interface{} {
	rval := reflect.ValueOf(data)
	val := reflect.Indirect(rval)
	if val.Kind() == reflect.Slice {
		val = val.Elem()
	}
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fieldNum := val.NumField()
	var ia []interface{}
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup("dbid"); ok {
			if tagval == "AUTOINCREMENT" || tagval == "SEQUENCE" {
				continue
			}
		}
		//only add parameters tagged with a db and value that is not empty or "-"
		if tagval, ok := typ.Field(i).Tag.Lookup("db"); ok {
			if tagval != "" && tagval != "-" {
				v := val.Field(i)
				if v.Kind() == reflect.Pointer && v.IsNil() {
					ia = append(ia, nil)
				} else {
					ia = append(ia, reflect.Indirect(v).Interface())
				}
			}
		}
		//if field is a struct, recursively traverse it.
		//most useful for encapulation
		if typ.Field(i).Type.Kind() == reflect.Struct {
			v := val.Field(i)
			embeddedParams := StructToIArray(reflect.Indirect(v).Interface())
			ia = append(ia, embeddedParams...)
		}
	}
	return ia
}

func StructToIArrayEx(data interface{}, excludeFields []string, tagField string, excludeTags []string) []interface{} {
	val := reflect.ValueOf(data).Elem()
	valtype := reflect.TypeOf(data).Elem()
	fieldNum := val.NumField()
	var ia []interface{}
	for i := 0; i < fieldNum; i++ {
		valField := val.Field(i)
		if excludeFields != nil {
			if contains(excludeFields, valtype.Field(i).Name) {
				continue
			}
		}
		if tagval, ok := valtype.Field(i).Tag.Lookup(tagField); ok {
			if contains(excludeTags, tagval) {
				continue
			}
		}
		ia = append(ia, valField.Addr().Interface())
	}
	return ia
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
