package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

var testYaml = `
one:
  test1:
    key: value
`

var defaultYaml = `
one:
  new: thing
  test1:
    key: value
two:
  test2:
    key: false
`

var mergeYaml = `
new: new_thing
`

func main() {
	values := Values{}
	merge := make(map[string]interface{})
	yaml.Unmarshal([]byte(testYaml), &values.customValues)
	yaml.Unmarshal([]byte(defaultYaml), &values.defaultValues)
	yaml.Unmarshal([]byte(mergeYaml), &merge)

	values.downLevel("one")
	fmt.Println(values.inCustomValues("test1"))
	values.upLevel()
	fmt.Println(values.inCustomValues("one"))
	/*	fmt.Println(values.inCustomValues("new"))
		values.downLevel("test1")
		values.upLevel()
		fmt.Println(values.currentCustomMap())
		values.mergeValues(merge)
		fmt.Println(values.inCustomValues("new"))
		values.downLevel("test1")
		fmt.Println(values.customValues)
		fmt.Println(values.currentCustomMap())
		fmt.Println(values.mergeSlice[len(values.mergeSlice)-1])
		fmt.Println(values.defaultValues)
		fmt.Println(values.currentKeys())
	*/
}
