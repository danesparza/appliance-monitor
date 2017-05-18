package data

import (
	"time"

	"github.com/boltdb/bolt"
)

// ConfigDB is the BoltDB database for config information
type ConfigDB struct {
	Database string
	User     string
	Password string
}

// ConfigItem represents a configuration item
type ConfigItem struct {
	ID          int64     `sql:"id" json:"id"`
	Name        string    `sql:"name" json:"name"`
	Value       string    `sql:"value" json:"value"`
	LastUpdated time.Time `sql:"updated" json:"updated"`
}

// InitStore initializes the database
func (store ConfigDB) InitStore(overwrite bool) error {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()

	return err
}
