# goquery

goquery is a small library to simplify commodity database operations. It provides a fluent API for database interactions with support for multiple database backends through a unified interface.

## Features

- Fluent API for database operations
- Support for PostgreSQL (via pgx) and other SQL databases (via sqlx)
- Type-safe data mapping with struct tags
- Transaction support
- JSON and CSV output formats
- Batch operations
- Connection pooling configuration

## Getting Started

### Installation

```bash
go get github.com/usace/goquery
```

For Postgres databases it wraps pgx and for sql interface db connections it wraps sqlx

### Basic Usage

```go
// Configure database connection
config := dq.RdbmsConfig{
    Dbuser:   "myuser",
    Dbpass:   "mypass",
    Dbhost:   "localhost",
    Dbname:   "postgres",
    Dbport:   "5432",
    DbDriver: "pgx",
    DbStore:  "pgx",
}

// Create data store
store, err := dq.NewRdbmsDataStore(&config)
if err != nil {
    log.Fatal(err)
}
```

## Core Concepts

### DataStore

The `DataStore` interface is the main entry point for database operations:

```go
type DataStore interface {
    Select(stmt string, params ...interface{}) *FluentSelect
    Get(dest interface{}, tx *Tx, stmt string, params ...interface{}) error
    Select(dest interface{}, tx *Tx, stmt string, params ...interface{}) error
    Query(tx *Tx, stmt string, params ...interface{}) (Rows, error)
    Insert(ds DataSet, rec interface{}, tx *Tx) error
    InsertStmt(ds DataSet) (string, error)
    Exec(tx *Tx, stmt string, params ...interface{}) error
    Execr(tx *Tx, stmt string, params ...interface{}) (ExecResult, error)
    MustExec(tx *Tx, stmt string, params ...interface{})
    MustExecr(tx *Tx, stmt string, params ...interface{}) ExecResult
    Transaction(fn func(tx Tx)) error
    Batch() (Batch, error)
    SendBatch(batch Batch) BatchResult
}
```

# Configuration

### Environment Variables

goquery supports configuration via environment variables:

- `DBUSER` - Database username
- `DBPASS` - Database password
- `DBHOST` - Database host
- `DBPORT` - Database port (default: 5432)
- `DBNAME` - Database name
- `DBDRIVER` - Database driver (e.g., "pgx")
- `DBSTORE` - Store type ("pgx" or "sqlx")
- `EXTERNAL_LIB` - External libraries for database connections
- `DBDRIVER_PARAMS` - Additional driver parameters
- `POOLMAXCONNS` - Maximum pool connections
- `POOLMINCONNS` - Minimum pool connections
- `POOLMAXCONNLIFETIME` - Maximum connection lifetime
- `POOLMAXCONNIDLE` - Maximum connection idle time

### RdbmsConfig Struct

```go
type RdbmsConfig struct {
    Dbuser      string
    Dbpass      string
    Dbhost      string
    Dbport      string
    Dbname      string
    ExternalLib string
    OnInit      string
    DbDriver    string
    DbStore     string

    PoolMaxConns        int
    PoolMinConns        int
    PoolMaxConnLifetime string //duration string
    PoolMaxConnIdle     string //duration string

    DbDriverSettings string
}
```

## Supported Databases

goquery supports multiple database backends:

- PostgreSQL (via pgx)
- Any SQL-compliant database (via sqlx)

## Database Dialects

The library includes support for different database dialects:

- PostgreSQL (`dialect_pg.go`)
- SQLite (`dialect_sqlite.go`)
- Oracle (`dialect_oracle.go`)
- DuckDB (`dialect_duckdb.go`)

Each dialect defines database-specific SQL syntax and connection methods.

## License

MIT License

---
 Creating a connection

 ```go
 store,err:=NewRdbmsDataStore(&config)
 ```

<br/>

---

### DataSet

DataSets are data structures used to define the structure of your data and organize associated SQL statements:

```go
type FishingSpot struct {
    ID       int32   `db:"id" dbid:"SEQUENCE" idsequence:"fishing_spots_id_seq"`
    Location *string `db:"location"`
}

var fs dq.TableDataSet = dq.TableDataSet{
    Name: "fishing_spots",
    Statements: dq.Statements{
        "get-fishing-spots":            "select * from fishing_spots",
        "get-fishing-spot-by-id":       "select * from fishing_spots where id=$1",
        "get-fishing-spot-by-location": "select * from fishing_spots where location=$1",
        "insert-with-return":           "insert into fishing_spots (location) values ($1) returning id",
    },
    TableFields: FishingSpot{}, //TableFields are only necessary for implciti insert statements
}
```

- Kitchen sink examples
```go

func postgresTest() {
	store, err := pgconnect()
	if err != nil {
		log.Fatalln(err)
	}

	////////////////////////////////////
	// Simple Select using in-line SQL
	// results are written to a slice
	////////////////////////////////////
	spots := []FishingSpot{}
	err = store.Select("select * from fishing_spots").
		Dest(&spots).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(spots)

	///////////////////////////////////
	// Simple select using a dataset reference
	// results are written to a slice
	//////////////////////////////////
	spots = []FishingSpot{}
	err = store.Select().
		DataSet(&fs).
		Dest(&spots).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(spots)

	///////////////////////////////////
	// Simple select using a dataset reference and a statement key
	// also includes parameter binding.
	// Params function takes a comma separated list of parameters of any type
	// that get bound to the query by position
	// results are written to a struct
	//////////////////////////////////
	dest := FishingSpot{}
	err = store.Select().
		DataSet(&fs).
		StatementKey("get-fishing-spot-by-location").
		Params("Pine Island").
		Dest(&dest).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(dest)

	///////////////////////////////////
	// Simple select that includes string concatonation via the Appliy function(yes...you read that properly)
	// SQL statement that is send to the DB is logged
	// results are written to a struct
	//////////////////////////////////
	dest = FishingSpot{}
	err = store.Select("select * from fishing_spots where %s=%d").
		Apply("id", 1).
		Dest(&dest).
		LogSql(true).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(dest)

	///////////////////////////////////
	// Simple select that returns a Rows cursor that the user can enumerate
	// Note that the caller is responsible for closing the Rows Cursor
	//////////////////////////////////
	rows, err := store.Select("select * from fishing_spots").
		FetchRows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		fs := FishingSpot{}
		rows.ScanStruct(&fs)
		columns, err := rows.Columns()
		if err != nil {
			log.Fatalln("failed to get column metadata")
		}

		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			log.Fatalln(err)
		}
		for i, col := range columns {
			switch col {
			case "id":
				fmt.Printf("col %s is %s and has a value of %v\n", col, columnTypes[i].Name(), fs.ID)
			case "location":
				fmt.Printf("col %s is %s and has a value of %s\n", col, columnTypes[i].Name(), safeprint(fs.Location))
			}
		}
	}

	///////////////////////////////////
	// Better approach for enumerating rows is to use the ForEachRow iterator
	// ForEachRow handles resource management and takes a function for the processing
	// of every row
	//////////////////////////////////
	dest = FishingSpot{}
	err = store.Select().
		DataSet(&fs).
		StatementKey("get-fishing-spots").
		ForEachRow(func(row dq.Rows) error {
			err := row.ScanStruct(&dest)
			if err != nil {
				return err
			}
			if dest.Location != nil {
				fmt.Printf("ROW FUNC:::%d:%s\n", dest.ID, *dest.Location)
			}
			return nil
		}).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}

	///////////////////////////////////
	// Select statements can apply a Suffix to a query or query statement key
	// The suffix will simply be added to the end of the query with a space separator
	//////////////////////////////////
	dest = FishingSpot{}
	err = store.Select("select * from fishing_spots").
		Suffix("where id=$1").
		Params(3).
		Dest(&dest).
		LogSql(true).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(dest)

	///////////////////////////////////
	// Fetch data as a json string
	// Using the FetchJsonMethod will accumunlate JSON bytes into a
	// byte buffer, so it should only be used for small datasets
	//////////////////////////////////
	json, err := store.Select("select * from fishing_spots").FetchJSON()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(json))

	///////////////////////////////////
	// More recent approach for fetching JSON data
	// OutputJson method takes a writer and the recordset cursor will
	// written directly to the writer
	//////////////////////////////////
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = store.Select("select * from fishing_spots").OutputJson(writer).Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(json))

	///////////////////////////////////
	// insert a single record without specifying an insert statement
	// in this case an insert statement will be generated and cached in the table dataset
	// with a statement key of "insert"
	//////////////////////////////////
	location := "100 location"
	spot := FishingSpot{
		ID:       100,
		Location: &location,
	}
	err = store.Insert(&fs).Records(&spot).Execute()

	data, err := store.Select("select * from fishing_spots where id=100").FetchJSON()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))

	///////////////////////////////////
	//remove the inserted record
	//uses query execution in the same manner as database/sql
	//////////////////////////////////
	err = store.Exec(dq.NoTx, "delete from fishing_spots where id=$1", 100)
	if err != nil {
		log.Fatalln(err)
	}

	///////////////////////////////////
	// insert a multiple record without specifying an insert statement
	// this will implicitly use the previously generated statement key of "insert"
	//////////////////////////////////

	location1 := "location1"
	location2 := "location2"

	locations := []FishingSpot{
		{
			ID:       101,
			Location: &location1,
		},
		{
			ID:       102,
			Location: &location2,
		},
	}

	err = store.Insert(&fs).Records(&locations).Execute()
	if err != nil {
		log.Fatalln(err)
	}

	data, err = store.Select("select * from fishing_spots where id>100").FetchJSON()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))

	///////////////////////////////////
	//remove the inserted records
	//uses query execution in the same manner as database/sql
	//////////////////////////////////
	err = store.Exec(dq.NoTx, "delete from fishing_spots where id>$1", 100)
	if err != nil {
		log.Fatalln(err)
	}

	///////////////////////////////////
	//insert multiple items within a transaction
	//statments must panic on failure inside transaction blocks
	//the transaction function will automatically rollback on error
	//and will atuomatically commit if there are no errors
	//////////////////////////////////

	err = store.Transaction(func(tx dq.Tx) {
		//use the insert fluent api
		store.Insert(&fs).Records(&locations[0]).Tx(&tx).PanicOnErr(true).Execute()

		//use a mustexec which will panic on error
		store.MustExec(&tx, "insert into fishing_spots values ($1,$2)", locations[1].ID, locations[1].Location)
	})

	data, err = store.Select("select * from fishing_spots where id>100").FetchJSON()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))

	///////////////////////////////////
	//remove the inserted records
	//but return the deleted ids.
	//anything with return values is effectively a select
	//this is similar to how you execute upserts
	//////////////////////////////////
	ids := []int{}
	err = store.Select(`
		DELETE FROM fishing_spots
		  WHERE id >100
		RETURNING id`).
		Dest(&ids).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ids)

	///////////////////////////////////
	//use the same appoach to perform an insert
	//but returning new id values
	//////////////////////////////////

	//set the fishing spots sequence to a value of 1000 since we were manually adding ids
	err = store.Exec(dq.NoTx, "SELECT setval('fishing_spots_id_seq', 1000, False)")
	if err != nil {
		log.Fatalln(err)
	}

	//execute the insert
	var newid int
	err = store.Select().
		DataSet(&fs).
		StatementKey("insert-with-return").
		Params("The new location").
		Dest(&newid).
		Fetch()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(newid)

}

func safeprint(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func pgconnect() (dq.DataStore, error) {
	config := dq.RdbmsConfig{
		Dbuser:   "myuser",
		Dbpass:   "mypass",
		Dbhost:   "localhost",
		Dbname:   "postgres",
		Dbport:   "5432",
		DbDriver: "pgx",
		DbStore:  "pgx",
	}
	return dq.NewRdbmsDataStore(&config)
}

```