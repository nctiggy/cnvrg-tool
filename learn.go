package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Contact struct {
	firstName   string
	lastName    string
	email       string
	phoneNumber string
	state       string
	business    bool
}

var contacts []Contact

var pages = tview.NewPages()

var app = tview.NewApplication()

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contact \n(q) to quit")

var form = tview.NewForm()

func addContactForm() {
	states := []string{"WA", "CA", "OR"}
	contact := Contact{}

	form.AddInputField("First Name", "", 20, nil, func(firstName string) {
		contact.firstName = firstName
	})

	form.AddInputField("Last Name", "", 20, nil, func(lastName string) {
		contact.lastName = lastName
	})

	form.AddInputField("Email", "", 20, nil, func(email string) {
		contact.email = email
	})

	form.AddInputField("Phone", "", 20, nil, func(phone string) {
		contact.phoneNumber = phone
	})

	// states is a slice of state abbreviations. Code is in the repo.
	form.AddDropDown("State", states, 0, func(state string, index int) {
		contact.state = state
	})

	form.AddCheckbox("Business", false, func(business bool) {
		contact.business = business
	})

	form.AddButton("Save", func() {
		contacts = append(contacts, contact)
		form.Clear(true)
		pages.SwitchToPage("Menu")
	})
}

func main() {
	pages.AddPage("Menu", text, true, true)
	pages.AddPage("Add Contact", form, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		text.SetText(string(event.Rune()))
		/*if event.Rune() == 27 {
			app.Stop()
		} else if event.Rune() == 'a' {
			addContactForm()
			pages.SwitchToPage("Add Contact")
		}*/
		return event
	})
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
