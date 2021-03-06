package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

const (
	configPrefix = "settings"
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

// Device represents a device in the system.
type Device struct {
	// ID is the unique device identifier
	ID string `json:"id"`

	// Type is the type of device (amon / hs110 / etc)
	Type string `json:"type"`

	// Name is the name given to the device by the customer
	Name string `json:"name"`

	// MinimumMonitorTime is the amount of time required to monitor the device before
	// it latches to the 'on' or 'off' state
	MinimumMonitorTime time.Duration `json:"minimum_monitor_time"`

	// Threshold is the value that must be crossed (after the MinimumMonitorTime time)
	// before the device latches to the 'on' or 'off' state
	Threshold int `json:"monitor_threshold"`

	// IPAddress is the network address for the device
	IPAddress string `json:"ipaddress"`

	// Running indicates whether the device is currently running (operating) or not
	Running bool `json:"running"`

	// LastUpdated indicates when this device was last updated in the system
	LastUpdated time.Time `sql:"updated" json:"updated"`
}

// InitStore initializes the database
func (store ConfigDB) InitStore() error {
	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()

	return err
}

// GetAllDevices gets all devices in the system
func (store ConfigDB) GetAllDevices() ([]Device, error) {
	//	Our return item
	retval := []Device{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all the items:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("devices"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			//	Unmarshal data into our config item
			device := Device{}
			if err := json.Unmarshal(v, &device); err != nil {
				return err
			}

			retval = append(retval, device)
		}

		return nil
	})

	//	Return our slice:
	return retval, err
}

// AddOrUpdateDevice adds a device to the system
func (store ConfigDB) AddOrUpdateDevice(device Device) (Device, error) {

	//	Our return item:
	retval := device

	//	If there is no config name, throw an error:
	if strings.TrimSpace(retval.Name) == "" {
		return retval, errors.New("Device name can't be blank")
	}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("devices"))
		if err != nil {
			return err
		}

		// If we don't have an id, generate an id for the configitem.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		if retval.ID == "" {
			retval.ID = xid.New().String()
		}

		//	Set the current datetime:
		retval.LastUpdated = time.Now()

		//	Serialize to JSON format
		encoded, err := json.Marshal(retval)
		if err != nil {
			return err
		}

		//	Store it, with the 'id' as the key:
		return b.Put([]byte(retval.ID), encoded)
	})

	return retval, err
}

// Set inserts or updates the config item
func (store ConfigDB) Set(configItem ConfigItem) (ConfigItem, error) {

	//	Our return item:
	retval := configItem

	//	If there is no config name, throw an error:
	if strings.TrimSpace(retval.Name) == "" {
		return retval, errors.New("Config name can't be blank")
	}

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
		if retval.ID == 0 {
			id, _ := b.NextSequence()
			retval.ID = int64(id)
		}

		//	Set the current datetime:
		retval.LastUpdated = time.Now()

		//	Serialize to JSON format
		encoded, err := json.Marshal(retval)
		if err != nil {
			return err
		}

		//	Store it, with the 'name' as the key:
		return b.Put([]byte(retval.Name), encoded)
	})

	return retval, err
}

// Remove removes the config item
func (store ConfigDB) Remove(configName string) error {

	//	If there is no config name, throw an error:
	if strings.TrimSpace(configName) == "" {
		return errors.New("Config name can't be blank")
	}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("configItems"))
		if err != nil {
			return err
		}

		//	Remove the item:
		return b.Delete([]byte(configName))
	})

	return err
}

// Get fetches a config item
func (store ConfigDB) Get(configName string) (ConfigItem, error) {
	//	Our return item:
	retval := ConfigItem{}

	//	Get the default from config file
	viperFormattedName := fmt.Sprintf("%v.%v", configPrefix, configName)
	viperConfigValue := viper.GetString(viperFormattedName)
	if viperConfigValue != "" {
		retval.Name = configName
		retval.Value = viperConfigValue
	}

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

			//	Need to make sure we got something back here before we try to unmarshal
			if len(configBytes) > 0 {
				//	Unmarshal data into our config item
				if err := json.Unmarshal(configBytes, &retval); err != nil {
					return err
				}
			}
		}

		return nil
	})

	//	If we found an item, use that .. otherwise, use the default
	return retval, err
}

// GetAll gets all config items in the system
func (store ConfigDB) GetAll() ([]ConfigItem, error) {
	//	Our return item:
	retval := []ConfigItem{}

	//	Get the defaults from config file
	if viper.InConfig(configPrefix) {
		items := viper.GetStringMapString(configPrefix)

		for k, v := range items {
			configItem := ConfigItem{
				Name:  k,
				Value: v}

			//	Add to the return slice:
			retval = append(retval, configItem)
		}
	}

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

			//	See if we can find it in our defaults.  If we can, just
			//	update the existing item.  Otherwise, add it
			foundConfigItem := false

			for i, v := range retval {
				if v.Name == configItem.Name {

					//	Update the item:
					retval[i].ID = configItem.ID
					retval[i].Value = configItem.Value
					retval[i].LastUpdated = configItem.LastUpdated

					//	Signal that we found the item:
					foundConfigItem = true
					break
				}
			}

			if !foundConfigItem {
				//	Add to the return slice:
				retval = append(retval, configItem)
			}

		}

		return nil
	})

	//	Return our slice:
	return retval, err
}
