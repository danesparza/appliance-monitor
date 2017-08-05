package data

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// EventType represents activity event types in the system
type EventType int

const (
	// ApplianceUknownState event is used when the state can't be determined
	ApplianceUknownState EventType = iota

	// ApplianceRunning event is signaled when the appliance appears to be running (vibrating)
	ApplianceRunning

	// ApplianceStopped event is signaled when the appliance appears to be stopped (not vibrating)
	ApplianceStopped

	// AppStarted event is signaled when the application starts
	AppStarted
)

// Activity represents a single activity event
type Activity struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"eventtype"`
}

// CloudActivity represents a single activity event in the cloud
type CloudActivity struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"eventtype"`
	DeviceID  string    `json:"deviceId"`
}

// ActivityDB is the BoltDB database for activity information
type ActivityDB struct {
	Database string
}

// Add inserts or updates activities
func (store ActivityDB) Add(activityItem Activity) (Activity, error) {

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

		//	Set the current datetime if needed:
		if activityItem.Timestamp.IsZero() {
			activityItem.Timestamp = time.Now()
		}

		//	Serialize to JSON format
		encoded, err := json.Marshal(activityItem)
		if err != nil {
			return err
		}

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(activityItem.Timestamp.Format(time.RFC3339)), encoded)
	})

	//	Set our return item:
	retval = activityItem

	return retval, err
}

// GetRange gets all activities in a given range
func (store ActivityDB) GetRange(startDate, endDate time.Time) ([]Activity, error) {
	retval := []Activity{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get the items in the given range:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("activities"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		// Format our timespan:
		min := []byte(startDate.Format(time.RFC3339))
		max := []byte(endDate.Format(time.RFC3339))

		// Iterate over the timespan
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

			//	Unmarshal data into our config item
			activity := Activity{}
			if err := json.Unmarshal(v, &activity); err != nil {
				return err
			}

			//	Add to the return slice:
			retval = append(retval, activity)
		}

		return nil
	})

	//	Return our slice:
	return retval, err
}

// GetAllActivity returns all activity
func (store ActivityDB) GetAllActivity() ([]Activity, error) {
	retval := []Activity{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all the items:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("activities"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			activity := Activity{}
			if err := json.Unmarshal(v, &activity); err != nil {
				return err
			}

			//	Add to the return slice:
			retval = append(retval, activity)
		}

		return nil
	})

	//	Return our slice:
	return retval, nil
}

// GetLatestActivity returns the most recent activity (or an empty Activity if no activity found)
func (store ActivityDB) GetLatestActivity() (Activity, error) {
	retval := Activity{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all the items:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("activities"))
		if b == nil {
			return nil
		}

		//	Get the last item:
		_, v := b.Cursor().Last()

		//	Unmarshal data into our config item
		if err := json.Unmarshal(v, &retval); err != nil {
			return err
		}

		return nil
	})

	//	Return our slice:
	return retval, nil
}

// DeleteRange removes all activities in a given range
func (store ActivityDB) DeleteRange(startDate, endDate time.Time) error {

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	//	Get the items in the given range:
	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("activities"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		// Format our timespan:
		min := []byte(startDate.Format(time.RFC3339))
		max := []byte(endDate.Format(time.RFC3339))

		// Iterate over the timespan
		for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {

			//	Delete the key:
			delerr := b.Delete(k)

			//	If we have an error, return early and with the error
			if delerr != nil {
				return err
			}
		}

		return nil
	})

	//	Return
	return nil
}
