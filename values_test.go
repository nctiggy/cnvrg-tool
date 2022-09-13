package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
)

var goodDefaultValues string = `
one:
  test1:
    key1: value1
    key2: 1
    key3: true
    key4:
      more:
        items: here
two:
  test2:
    key2: value2
`

var goodCustomValues string = `
one:
  test1:
    keyTest: valueTest
    key1: newValue
`

func buildGoodValues(v *Values) *Values {
	err := yaml.Unmarshal([]byte(goodDefaultValues), &v.defaultValues)
	if err != nil {
		fmt.Println("error loading default values")
	}
	yaml.Unmarshal([]byte(goodCustomValues), &v.customValues)
	if err != nil {
		fmt.Println("error loading current values")
	}
	return v
}

func TestDownLevel(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	goodValues.downLevel("one")
	parentLength := len(goodValues.parents)
	sliceLength := len(goodValues.mergeSlice)
	if parentLength != 1 {
		t.Errorf("Expected values.Parents length of 1 but got %d", parentLength)
	}
	if sliceLength != 1 {
		t.Errorf("Expected values.mergeSlice length of 1 but got %d", sliceLength)
	}

}

func TestUpLevel(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	goodValues.downLevel("two").upLevel()
	parentLength := len(goodValues.parents)
	sliceLength := len(goodValues.mergeSlice)
	if parentLength != 0 {
		t.Errorf("Expected values.Parents length of 1 but got %d", parentLength)
	}
	if sliceLength != 0 {
		t.Errorf("Expected values.mergeSlice length of 1 but got %d", sliceLength)
	}

}

func TestAddCustomLevel(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	goodValues.parents = append(goodValues.parents, "one")
	goodValues.addCustomLevel()
	if !reflect.DeepEqual(goodValues.mergeSlice[0], goodValues.customValues["one"]) {
		t.Errorf("Expected mergeSlice to match %v value", goodValues.customValues["one"])
	}
}

func TestRemoveCustomLevel(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	goodValues.parents = append(goodValues.parents, "one")
	goodValues.mergeSlice = append(goodValues.mergeSlice, goodValues.customValues["one"].(map[string]interface{}))
	goodValues.removeCustomLevel()
	if len(goodValues.parents) != 0 && len(goodValues.mergeSlice) != 0 {
		t.Errorf("Expected mergeSlice and parents to equal 0, instead equaled %d", len(goodValues.parents))
	}
}

func TestCurrentCustomMap(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	test := goodValues.currentCustomMap()
	if !reflect.DeepEqual(goodValues.customValues, test) {
		t.Errorf("Expected output (test) to match customValues attribute")
	}
	goodValues.downLevel("one")
	test1 := goodValues.currentCustomMap()
	if !reflect.DeepEqual(goodValues.customValues["one"], test1) {
		t.Errorf("Expected output (test1) to match customValues attribute")
	}
}

func TestCurrentDefaultMap(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	test := goodValues.currentDefaultMap()
	if !reflect.DeepEqual(goodValues.defaultValues, test) {
		t.Errorf("Expected output (test) to match defaultValues attribute")
	}
	goodValues.downLevel("one")
	test1 := goodValues.currentDefaultMap()
	if !reflect.DeepEqual(goodValues.defaultValues["one"], test1) {
		t.Errorf("Expected output (test1) to match defaultValues attribute")
	}
}

func TestCurrentKeys(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	testVal := make([]map[string]interface{}, 4)
	testVal[0] = map[string]interface{}{
		"name":     "key1",
		"changed":  '*',
		"hasChild": false,
		"value":    "newValue"}
	testVal[1] = map[string]interface{}{
		"name":     "key2",
		"changed":  ' ',
		"hasChild": false,
		"value":    1}
	testVal[2] = map[string]interface{}{
		"name":     "key3",
		"changed":  ' ',
		"hasChild": false,
		"value":    true}
	testVal[3] = map[string]interface{}{
		"name":     "key4",
		"changed":  ' ',
		"hasChild": true,
		"value":    nil}
	goodValues.downLevel("one").downLevel("test1")
	test := goodValues.currentKeys()
	sorted := sort.SliceIsSorted(test, func(i, j int) bool {
		return test[i]["name"].(string) < test[j]["name"].(string)
	})
	if !sorted {
		t.Errorf("Expected the slice to be alphabetized")
	}
	if len(test) != 4 {
		t.Errorf("Expeceted slice to have a length of 4, got %d", len(test))
	}
	for i := 0; i < len(test); i++ {
		if !reflect.DeepEqual(test[i], testVal[i]) {
			t.Errorf("Expected map to be %v was set to %v",
				testVal[i],
				test[i])
		}
	}
}

func TestIsDefaultValue(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	var badValue interface{} = "badValue"
	var goodValue interface{} = "value"
	correct := goodValues.downLevel("one").downLevel("test1").isDefaultValue("key1", goodValue)
	wrong := goodValues.isDefaultValue("key1", badValue)
	if !correct && wrong {
		t.Errorf("Default values are not correct")
	}
}

func TestInCustomValues(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	correct := goodValues.downLevel("one").downLevel("test1").inCustomValues("key1")
	wrong := goodValues.inCustomValues("key2")
	if !correct && wrong {
		t.Errorf("Custom values are not correct")
	}
}

func TestMergeValues(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	mergeMap := map[string]interface{}{"testKey": "testValue"}
	goodValues.downLevel("one").downLevel("test1").mergeValues(mergeMap)
	if !goodValues.inCustomValues("testKey") {
		t.Errorf("Merge did not happen correctly")
	}
}

func TestInit(t *testing.T) {
	goodValues := Values{}
	goodValues.init()
	if len(goodValues.defaultValues) != 0 && len(goodValues.customValues) != 0 {
		t.Errorf("expecting custom and default values to have a length of 0")
	}
}
