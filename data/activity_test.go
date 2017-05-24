package data_test

import (
	"os"
	"testing"

	"time"

	"github.com/danesparza/appliance-monitor/data"
)

//	Sanity check: The database shouldn't exist yet
func TestActivity_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("Activity database file check failed: Config file %s already exists, and shouldn't", filename)
	}
}

func TestActivity_Init_Successful(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ConfigDB{
		Database: filename}

	//	Act
	db.InitStore()

	//	Assert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Init failed: Activity db file %s was not created", filename)
	}
}

func TestActivity_Set_Successful(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	Try storing some config items:
	ct1 := data.Activity{
		DateTime: time.Now(),
		Type:     data.ApplianceRunning}

	//	Act
	response, err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: Should have set an item without error: %s", err)
	}

	if ct1.DateTime == response.DateTime {
		t.Errorf("Set failed: Should have set an item with the correct datetime: %+v / %+v", ct1, response)
	}
}
