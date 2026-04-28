package goquery

import (
	"database/sql/driver"
	"log"
	"os"
	"strconv"
)

// PoolMaxConnLifetime and PoolMaxConnIdle are string time duration representations
// as defined in ParseDuration in the stdlib time package
// the format consists of decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

const (
	dbPoolMaxConnsDefault int = 10
	dbPoolMinConnsDefault int = 2
)

type RdbmsConfig struct {
	Dbuser      string
	Dbpass      string
	Dbhost      string
	Dbport      string
	Dbname      string
	ExternalLib string
	OnInit      string //URL level initialization only supported by the Oracle GODROR driver
	DbDriver    string
	DbStore     string
	OnConnect   func(ds DataStore) error //optional function for running commands when a connection is opened

	PoolMaxConns        int
	PoolMinConns        int
	PoolMaxConnLifetime string //duration string
	PoolMaxConnIdle     string //duration string

	DbDriverSettings string

	//if a driver.connector is populated, it will be used to create all connections
	//all other parameters will be ignored
	Connector driver.Connector
}

func RdbmsConfigFromEnv() *RdbmsConfig {
	dbConfig := new(RdbmsConfig)
	dbConfig.Dbuser = os.Getenv("DBUSER")
	dbConfig.Dbpass = os.Getenv("DBPASS")
	dbConfig.Dbhost = os.Getenv("DBHOST")
	dbConfig.Dbport = os.Getenv("DBPORT")
	dbConfig.Dbname = os.Getenv("DBNAME")
	dbConfig.DbDriver = os.Getenv("DBDRIVER")
	dbConfig.DbStore = os.Getenv("DBSTORE")
	dbConfig.ExternalLib = os.Getenv("EXTERNAL_LIB")
	dbConfig.DbDriverSettings = os.Getenv("DBDRIVER_PARAMS")

	if dbConfig.Dbport == "" {
		dbConfig.Dbport = "5432"
	}

	maxConns := os.Getenv("POOLMAXCONNS")
	mc, err := strconv.Atoi(maxConns)
	if err != nil {
		log.Printf("Error parsing POOLMAXCONNS value of \"%s\":  Will fall back to default POOLMAXCONNS value of %d\n", maxConns, dbPoolMaxConnsDefault)
		dbConfig.PoolMaxConns = dbPoolMaxConnsDefault
	} else {
		dbConfig.PoolMaxConns = mc
	}

	minConns := os.Getenv("POOLMINCONNS")
	mc, err = strconv.Atoi(minConns)
	if err != nil {
		log.Printf("Error parsing POOLMINCONNS value of \"%s\":  Will fall back to default POOLMINCONNS value of %d\n", minConns, dbPoolMinConnsDefault)
		dbConfig.PoolMinConns = dbPoolMinConnsDefault
	} else {
		dbConfig.PoolMinConns = mc
	}

	dbConfig.PoolMaxConnLifetime = os.Getenv("POOLMAXCONNLIFETIME")
	dbConfig.PoolMaxConnIdle = os.Getenv("POOLMAXCONNIDLE")

	return dbConfig
}

/*
MaxConnLifetime time.Duration

	// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	MaxConnLifetimeJitter time.Duration

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime time.Duration

	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	MaxConns int32

	// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// to create new connections.
	MinConns int32

*/

/*
func (db *DB) SetConnMaxIdleTime(d time.Duration)
func (db *DB) SetConnMaxLifetime(d time.Duration)
func (db *DB) SetMaxIdleConns(n int)
func (db *DB) SetMaxOpenConns(n int)
*/
