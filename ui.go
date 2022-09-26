package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

var (
	app        = tview.NewApplication()
	pages      = tview.NewPages()
	list       = tview.NewList().ShowSecondaryText(false)
	form       = tview.NewForm()
	textView   = tview.NewTextView()
	hotKeyText = tview.NewTextView()
	flex       = tview.NewFlex()
	pad        = 1
)

func setText(s string, t *tview.TextView) {
	t.Clear()
	t.SetText(s)
}

func buildMenu(v *Values, k *hotKeys) {
	list.Clear()
	for _, key := range v.currentKeys() {
		var selFunc func() = nil
		name := key["name"].(string)
		if key["hasChild"].(bool) {
			text := name
			name = name + " >>"
			selFunc = func() {
				v.downLevel(text)
				buildMenu(v, k)
			}
		}
		list.AddItem(
			name,
			"",
			key["changed"].(rune),
			selFunc)
	}
	buildHotKeys(k, v)
	setText(k.generateHelp(), hotKeyText)
	k.setKeys(app)
}

func buildForm(v *Values, k *hotKeys) {
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
		buildMenu(v, k)
		pages.SwitchToPage("menu")
		app.SetFocus(list)
	})
}

func createHotKey(k tcell.Key, f func(), t string) hotKey {
	var key hotKey
	key.key = k
	key.action = f
	key.text = t
	return key
}

func buildHotKeys(k *hotKeys, v *Values) {
	k.resetKeys()
	hkEscape := createHotKey(
		tcell.KeyEscape,
		func() {
			v.upLevel()
			buildMenu(v, k)
		},
		"ESC: Menu Back")
	hkCtrlE := createHotKey(
		tcell.KeyCtrlE,
		func() {
			buildForm(v, k)
			pages.SwitchToPage("form")
			app.SetFocus(form)
		},
		"Ctrl-E: Edit Values")
	hkCtrlS := createHotKey(
		tcell.KeyCtrlS,
		func() {
			v.saveToFile("/tmp/", "cnvrg-values.yaml")
		},
		"Ctrl-S: Save Custom Values")
	k.addKey(hkCtrlE).addKey(hkEscape).addKey(hkCtrlS)
}

func main() {
	keys := hotKeys{}
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
	hotKeyText.SetText("default")
	flex.
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(pages, 0, 1, true).
			AddItem(textView, 0, 1, false), 0, 8, true).
		AddItem(hotKeyText, 0, 2, false)
	buildMenu(&values, &keys)
	/*
		cfg := config.GetConfigOrDie()
		discClient, _ := discovery.NewDiscoveryClientForConfig(cfg)
		vers, err1 := discClient.ServerVersion()
		if err1 != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		version := vers.String()
		setText("k8s ver: "+version, hotKeyText)
	*/
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
