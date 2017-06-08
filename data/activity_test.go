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
	response, err := db.Add(ct1)

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
	db.Add(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-5 * time.Minute),
		Type:      data.ApplianceStopped})

	//	Act
	response, err := db.GetRange(time.Now().Add(-4*time.Minute-10*time.Second), time.Now())

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
	db.Add(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceRunning})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceRunning})

	db.Add(data.Activity{
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

func TestActivity_DeleteRange_ItemsInRange_DeletesItems(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	Try storing some config items:
	db.Add(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-5 * time.Minute),
		Type:      data.ApplianceStopped})

	//	Act
	response1, err1 := db.GetRange(time.Now().Add(-10*time.Minute), time.Now())
	errDel := db.DeleteRange(time.Now().Add(-10*time.Minute), time.Now().Add(-4*time.Minute))
	response2, err2 := db.GetRange(time.Now().Add(-10*time.Minute), time.Now())

	//	Assert
	if err1 != nil {
		t.Errorf("Get range failed: Should have gotten the 1st range without error: %s", err1)
	}

	if err2 != nil {
		t.Errorf("Get range failed: Should have gotten the 2nd range without error: %s", err1)
	}

	if errDel != nil {
		t.Errorf("DeleteRange failed:  Should have removed the range without error: %s", errDel)
	}

	if len(response1) == len(response2) {
		t.Errorf("The remove didn't work: Should have gotten different counts.  Instead, got %v", len(response1))
	}

	if len(response2) != 3 {
		t.Errorf("The remove didn't remove the correct number: Should have 3 items left.  Instead, got %v", len(response2))
	}
}

func TestActivity_GetLatestActivity_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	No activities added

	//	Act
	response, err := db.GetLatestActivity()

	//	Assert
	if err != nil {
		t.Errorf("Get latest failed: Should have gotten the range item without error: %s", err)
	}

	if response.Type != data.ApplianceUknownState {
		t.Errorf("Get latest failed: Should have gotten default item back")
	}
}

func TestActivity_GetLatestActivity_ItemsInDB_ReturnsMostRecentItem(t *testing.T) {
	//	Arrange
	filename := "testactivity.db"
	defer os.Remove(filename)

	db := data.ActivityDB{
		Database: filename}

	//	Try storing some config items:
	db.Add(data.Activity{
		Timestamp: time.Now().Add(-1 * time.Minute),
		Type:      data.ApplianceRunning})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-2 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-3 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-4 * time.Minute),
		Type:      data.ApplianceStopped})

	db.Add(data.Activity{
		Timestamp: time.Now().Add(-5 * time.Minute),
		Type:      data.ApplianceStopped})

	//	Act
	response, err := db.GetLatestActivity()

	//	Assert
	if err != nil {
		t.Errorf("Get latest failed: Should have gotten the range item without error: %s", err)
	}

	if response.Type != data.ApplianceRunning {
		t.Errorf("Get latest failed: Should have returned most recent item.  Instead, got %v / %v", response.Timestamp, response.Type)
	}
}
