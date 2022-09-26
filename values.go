package main

import (
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type Values struct {
	defaultValues map[string]interface{}
	customValues  map[string]interface{}
	parents       []string
	mergeSlice    []map[string]interface{}
}

func (v *Values) init() *Values {
	v.defaultValues = make(map[string]interface{})
	v.customValues = make(map[string]interface{})
	return v
}

/*
func (v *Values) importDefaultValues(u string) *Values {

}
*/
func (v *Values) saveToFile(l string, f string) {
	data, _ := yaml.Marshal(v.customValues)
	os.WriteFile(l+f, data, 0644)
	setText("Saved to "+l+f, hotKeyText)
}

func (v *Values) mergeValues(m map[string]interface{}) *Values {
	var index int
	for i, _ := range v.parents {
		index = len(v.parents) - (i + 1)
		if i == 0 {
			for k, val := range m {
				if !v.isDefaultValue(k, val) {
					v.mergeSlice[index][k] = val
				} else {
					delete(v.mergeSlice[index], k)
				}
			}
		}
		if index == 0 {
			if len(v.mergeSlice[index]) == 0 {
				delete(v.customValues, v.parents[index])
			} else {
				v.customValues[v.parents[index]] = v.mergeSlice[index]
			}
		} else {
			if len(v.mergeSlice[index]) == 0 {
				delete(v.mergeSlice[index-1], v.parents[index])
			} else {
				v.mergeSlice[index-1][v.parents[index]] = v.mergeSlice[index]
			}
		}
	}
	return v
}

func (v *Values) upLevel() *Values {
	if length := len(v.parents); length > 0 {
		v.parents = v.parents[:length-1]
		v.removeCustomLevel()
	}
	return v
}

func (v *Values) downLevel(k string) *Values {
	v.parents = append(v.parents, k)
	v.addCustomLevel()
	return v
}

func (v *Values) addCustomLevel() {
	currentLevel := make(map[string]interface{})
	for k, val := range v.customValues {
		currentLevel[k] = val
	}
	for _, parent := range v.parents {
		if _, ok := currentLevel[parent]; ok {
			currentLevel = currentLevel[parent].(map[string]interface{})
		} else {
			currentLevel = make(map[string]interface{})
			break
		}
	}
	v.mergeSlice = append(v.mergeSlice, currentLevel)
}

func (v *Values) removeCustomLevel() {
	l := len(v.mergeSlice)
	v.mergeSlice = v.mergeSlice[:l-1]
}

func (v *Values) currentDefaultMap() map[string]interface{} {
	currentLevel := make(map[string]interface{})
	for k, val := range v.defaultValues {
		currentLevel[k] = val
	}
	for _, parent := range v.parents {
		currentLevel = currentLevel[parent].(map[string]interface{})
	}
	return currentLevel
}

func (v *Values) currentCustomMap() map[string]interface{} {
	if len(v.parents) == 0 {
		return v.customValues
	} else {
		return v.mergeSlice[len(v.parents)-1]
	}
}

func (v *Values) currentKeys() []map[string]interface{} {
	keys := []map[string]interface{}{}
	currentMap := v.currentDefaultMap()
	currentCustomMap := v.currentCustomMap()
	for key, value := range currentMap {
		tempMap := make(map[string]interface{})
		tempMap["name"] = key
		tempMap["changed"] = ' '
		setVal := value
		if val, ok := currentCustomMap[key]; ok {
			setVal = val
		}
		if v.inCustomValues(key) {
			tempMap["changed"] = '*'
		}
		switch value.(type) {
		case map[string]interface{}:
			if len(value.(map[string]interface{})) > 0 {
				tempMap["hasChild"] = true
				tempMap["value"] = nil
			} else {
				tempMap["hasChild"] = false
				tempMap["value"] = setVal
			}
		default:
			tempMap["hasChild"] = false
			tempMap["value"] = setVal
		}
		keys = append(keys, tempMap)
	}
	sort.Slice(
		keys,
		func(i, j int) bool {
			return keys[i]["name"].(string) < keys[j]["name"].(string)
		})
	return keys
}

func (v *Values) isDefaultValue(key string, inVal interface{}) bool {
	defaultValues := v.currentDefaultMap()
	if value, ok := defaultValues[key]; ok {
		if inVal == value {
			return true
		}
		return false
	}
	return false
}

func (v *Values) inCustomValues(key string) bool {
	currentLevel := v.currentCustomMap()
	if _, ok := currentLevel[key]; ok {
		return true
	}
	return false
}
