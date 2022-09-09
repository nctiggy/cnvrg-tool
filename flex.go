package main

import (
	"github.com/rivo/tview"
)

var box tview.NewBox()

func main() {
	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("one", "", 'a', nil).
		AddItem("two", "", 'b', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	flex := tview.NewFlex().
		AddItem(list.SetBorder(true), 0, 1, false).
		AddItem(tview.NewBox().
			SetBorder(true).
			SetTitle("stuff here"), 0, 2, false)
	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
