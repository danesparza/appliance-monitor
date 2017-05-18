package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/appliance-monitor/data"
)

//	Sanity check: The database shouldn't exist yet
func TestConfig_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := "testing.db"

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("Config database file check failed: Config file %s already exists, and shouldn't", filename)
	}
}

//	Init should create a new BoltDB file
func TestConfig_Init_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := data.ConfigDB{
		Database: filename}

	//	Act
	db.InitStore()

	//	Assert
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Init failed: Config db file %s was not created", filename)
	}
}

//	Config get should return successfully even if the item doesn't exist
func TestConfig_Get_ItemDoesntExist_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := data.ConfigDB{
		Database: filename}

	queryName := "bogusItem"
	expectedValue := ""

	//	Act
	response, err := db.Get(queryName)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned an empty dataset without error: %s", err)
	}

	if expectedValue != response.Value && response.Value != "" {
		t.Errorf("Get failed: Shouldn't have returned the value %s", response.Value)
	}
}

//	Config set should work
func TestConfig_Set_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := data.ConfigDB{
		Database: filename}

	//	Try storing some config items:
	ct1 := data.ConfigItem{
		Name:  "TestItem1",
		Value: "Value1"}

	//	Act
	response, err := db.Set(ct1)

	//	Assert
	if err != nil {
		t.Errorf("Set failed: Should have set an item without error: %s", err)
	}

	if ct1.ID == response.ID {
		t.Errorf("Set failed: Should have set an item with the correct id: %+v / %+v", ct1, response)
	}
}

//	Config set then get should work
func TestConfig_Set_ThenGet_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)

	db := data.ConfigDB{
		Database: filename}

	//	Try storing some config items:
	ct1 := data.ConfigItem{
		Name:  "TestItem1",
		Value: "Value1"}

	ct2 := data.ConfigItem{
		Name:  "TestItem2",
		Value: "Value2"}

	queryName := "TestItem1"
	expectedValue := "Value1"

	//	Act
	db.Set(ct1)
	db.Set(ct2)
	response, err := db.Get(queryName)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned a config item without error: %s", err)
	}

	if response.Value != expectedValue {
		t.Errorf("Get failed: Should have returned the value %s but returned %s instead", ct2.Value, response.Value)
	}
}
