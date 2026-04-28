package goquery

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
)

func TestDuckDb_SpatialQuery(t *testing.T) {
	config2 := RdbmsConfig{
		DbDriver: "duckdb",
		DbStore:  "sqlx",
		OnConnect: func(db DataStore) error {
			return db.Exec(NoTx, "INSTALL sqlite; LOAD sqlite; INSTALL spatial; LOAD spatial")
		},
	}

	store2, err := NewRdbmsDataStore(&config2)
	if err != nil {
		log.Fatal(err)
	}

	resource := "testdata/example.gpkg"

	query := fmt.Sprintf("SELECT * FROM ST_Read('%s',layer='%s', keep_wkb=true)", resource, "geometry1")
	err = store2.Select(query).ForEachRow(func(row Rows) error {
		valsMap, err := row.ToMap()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(valsMap["text"])
		return nil
	}).Fetch()
}
