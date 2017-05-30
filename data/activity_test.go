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
		Timestamp: time.Now(),
		Type:      data.ApplianceRunning}

	//	Act
	response, err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: Should have set an item without error: %s", err)
	}

	if response.Timestamp.IsZero() {
		t.Errorf("Set failed: Should have set an item with the correct datetime: %+v", response)
	}
}

func TestActivity_GetRange_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	No items are in the database!

	//	Act
	response, err := db.GetRange(time.Now().Add(-10*time.Minute), time.Now())

	//	Assert
	if err != nil {
		t.Errorf("Get range failed: Should have gotten the range item without error: %s", err)
	}

	if len(response) != 0 {
		t.Errorf("Get range failed: Should not have gotten any items")
	}
}

func TestActivity_GetRange_ItemsInRange_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	Try storing some config items:
	db.Set(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-5 * time.Minute),
		Type:      data.ApplianceStopped})

	//	Act
	response, err := db.GetRange(time.Now().Add(-4*time.Minute), time.Now())

	//	Assert
	if err != nil {
		t.Errorf("Get range failed: Should have gotten the range item without error: %s", err)
	}

	if len(response) != 4 {
		t.Errorf("Get range failed: Should have gotten all items.  Instead, got %v", len(response))
	}
}

func TestActivity_GetAllActivity_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	Try storing some config items:
	db.Set(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceRunning})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceRunning})

	db.Set(data.Activity{
		Timestamp: time.Now().Add(-5 * time.Minute),
		Type:      data.ApplianceStopped})

	//	Act
	response, err := db.GetAllActivity()

	//	Assert
	if err != nil {
		t.Errorf("Get range failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 5 {
		t.Errorf("Get range failed: Should have gotten all items")
	}
}
