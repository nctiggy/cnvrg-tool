package main

import (
	"github.com/rivo/tview"
)

type pages struct {
	pages tview.Pages
}

type menu struct {
	list tview.List
}

func (m *menu) init() {
	m.list = tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(false)
}

func (m *menu) buildMenu(v *Values) {
	m.list.Clear()

}
