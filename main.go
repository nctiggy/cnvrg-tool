package main

import (
	"fmt"

	"github.com/rivo/tview"
	"gopkg.in/yaml.v2"
)

var data = `
- name: controlplane
  group: true
- name: networking
  group: true
- name: backup
  questions:
    enabled: bool
    retention: string
    interval: string
    httpProxy: slice
    ingressSvcAnnotations: map
  parent: controlplane
- name: capsule
  questions:
    image: string
  parent: networking
- name: labels
  type: map
  group: true
`

type Section struct {
	Name       string                 `yaml:"name"`
	Questions  map[string]interface{} `yaml:"questions",omitempty`
	Answers    map[string]interface{} `yaml:"answers",omitempty`
	Parent     string                 `yaml:"parent",omitempty`
	ParentTree []string               `yaml:"parentTree",omitempty`
	Type       string                 `yaml:"type",omitempty`
	Group      bool                   `yaml:"group",omitempty`
}

type m = map[string]interface{}

func mergeKeys(left, right m) m {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			//then we don't want to replace it - recurse
			left[key] = mergeKeys(leftVal.(m), rightVal.(m))
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}

func (s Section) hasParent() bool {
	if s.Parent != "" {
		return true
	}
	return false
}

func parentNames(s *Section, p *[]string, sections []Section) {
	*p = append(*p, s.Name)
	if s.hasParent() {
		for _, v := range sections {
			if v.Name == s.Parent {
				parentNames(&v, p, sections)
			}
		}
	}
}

func captureString(k string) string {
	answer := "test"
	return answer
}

func captureInt(k string) int {
	answer := 2
	return answer
}

func captureSlice(k string) []string {
	answer := []string{"test", "test1"}
	return answer
}

func captureMap(k string) m {
	answer := make(map[string]interface{})
	answer["testKey"] = "testValue"
	answer["testKey1"] = "testValue"
	answer["testKey2"] = "testValue"
	return answer
}

func captureBool(k string) bool {
	answer := true
	return answer
}

func genericQuestion(k string) {
	return
}

func (s Section) askQuestions() m {
	r := make(map[string]interface{})

	if s.hasEnabled() {
		fmt.Printf("Enable %v: \n", s.Name)
		r["enabled"] = true //fake answer
	}
	if v, ok := r["enabled"]; !ok || v == true {
		for k, v := range s.Questions {
			fmt.Printf("Set %v:\n", k)
			switch v {
			case "string":
				r[k] = captureString(k)
			case "bool":
				r[k] = captureBool(k)
			case "slice":
				r[k] = captureSlice(k)
			case "map":
				r[k] = captureMap(k)
			case "int":
				r[k] = captureInt(k)
			}
		}
	}
	return r
}

func (s *Section) hasEnabled() bool {
	if _, ok := s.Questions["enabled"]; ok {
		delete(s.Questions, "enabled")
		return true
	}
	return false
}

func addSection(p []string, a m, y m) {
	if len(a) == 0 {
		return
	}
	temp := make(map[string]interface{})
	for i, v := range p {
		dummy := make(map[string]interface{})
		switch {
		case i == 0:
			dummy[v] = a
		//case i == len(p)-1:
		//	continue
		default:
			dummy[v] = temp
		}
		temp = dummy
	}
	y = mergeKeys(y, temp)
}

func generateMenu(list *tview.List, s *Section) {
	entry := fmt.Sprintf("%s section", s.Name)
	list.AddItem(entry, "", 0, func() {
		for
	})
}

func main() {
	app := tview.NewApplication()
	output := make(map[string]interface{})
	var sections []Section
	err := yaml.Unmarshal([]byte(data), &sections)
	if err != nil {
		fmt.Println(err)
	}

	list := tview.NewList().ShowSecondaryText(false)
	flex := tview.NewFlex().
		AddItem(list, 0, 3, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("stuff here"), 0, 2, false)
	for _, section := range sections {
		if section.Group {
			generateMenu(list, &section)
		}
		answers := make(map[string]interface{})
		if section.Questions != nil {
			answers = section.askQuestions()
		} else if section.Type == "map" {
			answers = captureMap(section.Name)
		} else {
			continue
		}
		parents := []string{}
		parentNames(&section, &parents, sections)
		addSection(parents, answers, output)
	}
	d, err := yaml.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", string(d))
	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
