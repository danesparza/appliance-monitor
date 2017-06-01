package data

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// ConfigDB is the BoltDB database for config information
type ConfigDB struct {
	Database string
}

// ConfigItem represents a configuration item
type ConfigItem struct {
	ID          int64     `sql:"id" json:"id"`
	Name        string    `sql:"name" json:"name"`
	Value       string    `sql:"value" json:"value"`
	LastUpdated time.Time `sql:"updated" json:"updated"`
}

// InitStore initializes the database
func (store ConfigDB) InitStore() error {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()

	return err
}

// Set inserts or updates the config item
func (store ConfigDB) Set(configItem ConfigItem) (ConfigItem, error) {

	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("configItems"))
		if err != nil {
			return err
		}

		// If we don't have an id, generate an id for the configitem.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		if configItem.ID == 0 {
			id, _ := b.NextSequence()
			configItem.ID = int64(id)
		}

		//	Set the current datetime:
		configItem.LastUpdated = time.Now()

		//	Serialize to JSON format
		encoded, err := json.Marshal(configItem)
		if err != nil {
			return err
		}

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(configItem.Name), encoded)
	})

	//	Set our return item:
	retval = configItem

	return retval, err
}

// Get fetches a config item
func (store ConfigDB) Get(configName string) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		//	Get the item from the bucket
		b := tx.Bucket([]byte("configItems"))

		if b != nil {
			configBytes := b.Get([]byte(configName))

			//	Need to make sure we got something back here before we try to unmarshal?
			if len(configBytes) > 0 {
				//	Unmarshal data into our config item
				if err := json.Unmarshal(configBytes, &retval); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return retval, err
}

// GetAll gets all config items in the system
func (store ConfigDB) GetAll() ([]ConfigItem, error) {
	//	Our return item:
	retval := []ConfigItem{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all the items:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("configItems"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			configItem := ConfigItem{}
			if err := json.Unmarshal(v, &configItem); err != nil {
				return err
			}

			//	Add to the return slice:
			retval = append(retval, configItem)
		}

		return nil
	})

	//	Return our slice:
	return retval, err
}
