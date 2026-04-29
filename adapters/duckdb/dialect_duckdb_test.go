package duckdb

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/usace/goquery/v3"
)

func TestDuckDb_SpatialQuery(t *testing.T) {
	config2 := goquery.RdbmsConfig{
		DbDriver: "duckdb",
		DbStore:  "sqlx",
		OnConnect: func(db goquery.DataStore) error {
			return db.Exec(goquery.NoTx, "INSTALL sqlite; LOAD sqlite; INSTALL spatial; LOAD spatial")
		},
	}

	store2, err := goquery.NewRdbmsDataStore(&config2)
	if err != nil {
		log.Fatal(err)
	}

	resource := "testdata/example.gpkg"
	layer := "geometry1"

	query := fmt.Sprintf("SELECT * FROM ST_Read('%s',layer='%s', keep_wkb=true)", resource, layer)
	err = store2.Select(query).
		ForEachRow(func(row goquery.Rows) error {
			valsMap, err := row.ToMap()
			if err != nil {
				t.Error(err)
			}
			fmt.Println(valsMap["text"])
			return nil
		}).Fetch()
}
