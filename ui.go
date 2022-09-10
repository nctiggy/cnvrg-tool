package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

var app = tview.NewApplication()

var pages = tview.NewPages()

var list = tview.NewList().ShowSecondaryText(false)

var form = tview.NewForm()

var textView = tview.NewTextView()

var hotKeyText = tview.NewTextView()

var flex = tview.NewFlex()

var pad int = 1

func setText(s string, t *tview.TextView) {
	t.Clear()
	t.SetText(s)
}

func buildMenu(v *Values) {
	list.Clear()
	for _, key := range v.currentKeys() {
		var selFunc func() = nil
		name := key["name"].(string)
		if key["hasChild"].(bool) {
			text := name
			name = name + " >>"
			selFunc = func() {
				v.downLevel(text)
				buildMenu(v)
			}
		}
		list.AddItem(
			name,
			"",
			key["changed"].(rune),
			selFunc)
	}
	setEscape(v)
}

func buildForm(v *Values) {
	form.Clear(true)
	formEntries := make(map[string]interface{})
	for _, key := range v.currentKeys() {
		if key["hasChild"].(bool) {
			continue
		}
		switch key["value"].(type) {
		case string:
			form.AddInputField(
				key["name"].(string),
				key["value"].(string),
				20,
				nil,
				func(value string) {
					i, _ := form.GetFocusedItemIndex()
					formEntries[form.GetFormItem(i).GetLabel()] = value
				})
		case bool:
			form.AddCheckbox(
				key["name"].(string),
				key["value"].(bool),
				func(value bool) {
					i, _ := form.GetFocusedItemIndex()
					formEntries[form.GetFormItem(i).GetLabel()] = value
				})
		case int:
			form.AddInputField(
				key["name"].(string),
				strconv.Itoa(key["value"].(int)),
				20,
				nil,
				func(value string) {
					i, _ := form.GetFocusedItemIndex()
					formEntries[form.GetFormItem(i).GetLabel()], _ = strconv.Atoi(value)
				})
		case map[string]interface{}:
			form.AddInputField(
				key["name"].(string),
				"",
				20,
				nil,
				func(value string) {
					i, _ := form.GetFocusedItemIndex()
					formEntries[form.GetFormItem(i).GetLabel()] = value
				})
		}

	}
	form.AddButton("Cancel", func() {
		form.Clear(true)
		pages.SwitchToPage("menu")
		app.SetFocus(list)
	})
	form.AddButton("Save", func() {
		v.mergeValues(formEntries)
		d, _ := yaml.Marshal(v.customValues)
		setText(string(d), textView)
		form.Clear(true)
		buildMenu(v)
		pages.SwitchToPage("menu")
		app.SetFocus(list)
	})
}

func setEscape(v *Values) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			v.upLevel()
			buildMenu(v)
		case tcell.KeyCtrlE:
			buildForm(v)
			pages.SwitchToPage("form")
			app.SetFocus(form)
		}
		return event
	})
}

func main() {
	values := Values{}
	values.init()
	yamlFile, _ := ioutil.ReadFile("values.yaml")
	err := yaml.Unmarshal(yamlFile, &values.defaultValues)
	if err != nil {
		fmt.Println(err)
	}
	pages.
		SetBorder(true).
		SetTitle("Menu").
		SetBorderPadding(pad, pad, pad, pad)
	pages.
		AddPage("menu", list, true, true).
		AddPage("form", form, true, false)
	list.
		SetWrapAround(false)
	textView.
		SetTitle("values.yaml Preview").
		SetBorder(true).
		SetBorderPadding(pad, pad, pad, pad)
	hotKeyText.
		SetTitle("Hotkeys").
		SetBorder(true).
		SetBorderPadding(pad, pad, pad, pad)
	flex.
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(pages, 0, 1, true).
			AddItem(textView, 0, 1, false), 0, 9, true).
		AddItem(hotKeyText, 0, 1, false)
	buildMenu(&values)
	//for _, v := range defaultValues {
	//varType := fmt.Sprint(reflect.TypeOf(v))
	//}
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
