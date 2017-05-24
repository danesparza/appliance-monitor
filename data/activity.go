package data

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// EventType represents activity event types in the system
type EventType int

const (
	// ApplianceRunning event is signaled when the appliance appears to be running (vibrating)
	ApplianceRunning EventType = iota

	// ApplianceStopped event is signaled when the appliance appears to be stopped (not vibrating)
	ApplianceStopped

	// AppStarted event is signaled when the application starts
	AppStarted
)

// Activity represents a single activity event
type Activity struct {
	DateTime time.Time `json:"time"`
	Type     EventType `json:"eventtype"`
}

// ActivityDB is the BoltDB database for activity information
type ActivityDB struct {
	Database string
}

// Set inserts or updates the config item
func (store ActivityDB) Set(activityItem Activity) (Activity, error) {

	//	Our return item:
	retval := Activity{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("activities"))
		if err != nil {
			return err
		}

		//	Set the current datetime:
		activityItem.DateTime = time.Now()

		//	Serialize to JSON format
		encoded, err := json.Marshal(activityItem)
		if err != nil {
			return err
		}

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(activityItem.DateTime.Format(time.RFC3339)), encoded)
	})

	//	Set our return item:
	retval = activityItem

	return retval, err
}
