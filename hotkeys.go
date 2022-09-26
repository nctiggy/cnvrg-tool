package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type hotKey struct {
	key    tcell.Key
	action func()
	text   string
}

type hotKeys struct {
	keys []hotKey
}

func (k *hotKeys) resetKeys() *hotKeys {
	k.keys = k.keys[:0]
	return k
}

func (k *hotKeys) addKey(h hotKey) *hotKeys {
	k.keys = append(k.keys, h)
	return k
}

func (k *hotKeys) generateHelp() string {
	var help string
	for _, v := range k.keys {
		if help == "" {
			help = v.text
		} else {
			help = fmt.Sprintf("%s | %s", help, v.text)
		}
	}
	return help
}

func (k *hotKeys) setKeys(a *tview.Application) *hotKeys {
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		for _, v := range k.keys {
			if v.key == event.Key() {
				v.action()
			}
		}
		return event
	})
	return k
}
