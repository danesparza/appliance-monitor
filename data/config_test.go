package data_test

import (
	"bytes"
	"os"
	"testing"

	"fmt"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/spf13/viper"
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
	defer viper.Reset()

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
	defer viper.Reset()

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

//	Config get should return default even if the item doesn't exist in database
func TestConfig_Get_ItemDoesntExistButHasDefault_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	viper.SetConfigType("yaml")
	var yamlConfig = []byte(`
settings:
  afirstconfigitem: somethinghere
  itemwithdefault: thedefault
  somethingelse: anothervalue
`)
	viper.ReadConfig(bytes.NewBuffer(yamlConfig)) // Read in the defaults from the config file

	db := data.ConfigDB{
		Database: filename}

	queryName := "itemwithdefault"
	expectedValue := "thedefault"

	//	Act
	response, err := db.Get(queryName)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned the default without error: %s", err)
	}

	if expectedValue != response.Value {
		t.Errorf("Get failed: Should have returned the default '%v' instead of the value '%s'", expectedValue, response.Value)
	}
}

//	Config get should return default even if the item doesn't exist in database
func TestConfig_Get_ItemExistsAndHasDefault_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	viper.SetConfigType("yaml")
	var yamlConfig = []byte(`
settings:
  afirstconfigitem: somethinghere
  itemwithdefault: thedefault
  somethingelse: anothervalue
`)
	viper.ReadConfig(bytes.NewBuffer(yamlConfig)) // Read in the defaults from the config file

	db := data.ConfigDB{
		Database: filename}

	db.Set(data.ConfigItem{
		Name:  "itemwithdefault",
		Value: "newvalue"})

	queryName := "itemwithdefault"
	expectedValue := "newvalue"

	//	Act
	response, err := db.Get(queryName)

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned the default without error: %s", err)
	}

	if expectedValue != response.Value {
		t.Errorf("Get failed: Should have returned '%v' instead of the value '%s'", expectedValue, response.Value)
	}
}

//	Config set with no config name shouldn't work
func TestConfig_Set_NoName_NotSuccessful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	db := data.ConfigDB{
		Database: filename}

	//	Try storing some config items:
	ct1 := data.ConfigItem{
		Name:  "",
		Value: "Value1"}

	//	Act
	_, err := db.Set(ct1)

	//	Assert
	if err == nil {
		t.Errorf("Set failed: Should have thrown an error about no config name")
	}
}

//	Config set should work
func TestConfig_Set_ValidName_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

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

	if ct1.Name != response.Name {
		t.Errorf("Set failed: Should have set an item with the correct name: %+v / %+v", ct1, response)
	}

}

//	Config set then get should work
func TestConfig_Set_ThenGet_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

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
		t.Errorf("Set then Get failed: Should have returned a config item without error: %s", err)
	}

	if response.Value != expectedValue {
		t.Errorf("Set then Get failed: Should have returned the value '%s' but returned '%s' instead", ct2.Value, response.Value)
	}

	if response.Name != queryName {
		t.Errorf("Set then Get failed: Should have returned the value '%s' but returned '%s' instead", queryName, response.Name)
	}
}

func TestConfig_GetAll_NoItems_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	db := data.ConfigDB{
		Database: filename}

	//	NO ITEMS STORED

	//	Act
	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("GetAll failed: Should have returned config items without error: %s", err)
	}

	if len(response) != 0 {
		t.Errorf("GetAll failed: Should have returned no config items but returned %v instead", len(response))
	}
}

func TestConfig_GetAll_WithItems_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	db := data.ConfigDB{
		Database: filename}

	//	Try storing some config items:
	maxItems := 20
	for c := 1; c <= maxItems; c++ {
		db.Set(data.ConfigItem{
			Name:  fmt.Sprintf("TestItem%d", c),
			Value: fmt.Sprintf("Value %d", c)})
	}

	//	Act
	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("GetAll failed: Should have returned config items without error: %s", err)
	}

	if len(response) != maxItems {
		t.Errorf("GetAll failed: Should have returned %d config items but returned %v instead", maxItems, len(response))
	}
}

//	Config get should return default even if the item doesn't exist in database
func TestConfig_GetAll_NoItemsButHasDefaults_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	viper.SetConfigType("yaml")
	var yamlConfig = []byte(`
settings:
  afirstconfigitem: somethinghere
  itemwithdefault: thedefault
  somethingelse: anothervalue
`)
	viper.ReadConfig(bytes.NewBuffer(yamlConfig)) // Read in the defaults from the config file

	db := data.ConfigDB{
		Database: filename}

	expectedCount := 3

	//	Act
	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned the defaults without error: %s", err)
	}

	if expectedCount != len(response) {
		t.Errorf("GetAll failed: Should have returned the expected number of items %v instead of %v", expectedCount, len(response))
	}

}

//	Config get should return default even if the item doesn't exist in database
func TestConfig_GetAll_ItemsAndDefaults_Successful(t *testing.T) {
	//	Arrange
	filename := "testing.db"
	defer os.Remove(filename)
	defer viper.Reset()

	viper.SetConfigType("yaml")
	var yamlConfig = []byte(`
settings:
  afirstconfigitem: somethinghere
  itemwithdefault: thedefault
  somethingelse: anothervalue
`)
	viper.ReadConfig(bytes.NewBuffer(yamlConfig)) // Read in the defaults from the config file

	db := data.ConfigDB{
		Database: filename}

	db.Set(data.ConfigItem{
		Name:  "itemwithdefault",
		Value: "newvalue"})

	expectedCount := 3
	foundConfigItem := false
	configName := "itemwithdefault"
	expectedValue := "newvalue"

	//	Act
	response, err := db.GetAll()

	//	Assert
	if err != nil {
		t.Errorf("Get failed: Should have returned the defaults without error: %s", err)
	}

	if expectedCount != len(response) {
		t.Errorf("GetAll failed: Should have returned the expected number of items %v instead of %v", expectedCount, len(response))
	}

	for _, v := range response {
		if v.Name == configName {

			if v.Value != expectedValue {
				t.Errorf("GetAll failed: Should have returned the new value '%s' but got '%s'", expectedValue, v.Value)
			}

			//	Signal that we found the item:
			foundConfigItem = true
			break
		}
	}

	if !foundConfigItem {
		t.Errorf("GetAll failed: Couldn't find the config item '%s'", configName)
	}
}
