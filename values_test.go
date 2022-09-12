package main

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

var goodDefaultValues string = `
one:
  test1:
    key1: value1
two:
  test2:
    key2: value2
`

var goodCustomValues string = `
one:
  test1:
    keyTest: valueTest
`

func buildGoodValues(v *Values) *Values {
	yaml.Unmarshal([]byte(goodDefaultValues), &v.defaultValues)
	yaml.Unmarshal([]byte(goodCustomValues), &v.customValues)
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

func TestCurrentDefaultMap(t *testing.T) {
	goodValues := Values{}
	buildGoodValues(&goodValues)
	test := goodValues.currentCustomMap()
	if !reflect.DeepEqual(goodValues.customValues, test) {
		t.Errorf("Expected output to match customValues attribute")
	}
}
