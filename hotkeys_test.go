package main

import (
	"strconv"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func buildKeys(h *hotKeys) {
	keys := []tcell.Key{
		tcell.KeyEscape,
		tcell.KeyDelete,
		tcell.KeyEnter,
		tcell.KeyTab}
	for i := 0; i < 4; i++ {
		var key hotKey
		key.key = keys[i]
		key.action = func() {}
		key.text = "text" + strconv.Itoa(i)
		h.keys = append(h.keys, key)
	}
}

func TestResetKeys(t *testing.T) {
	var keys hotKeys
	buildKeys(&keys)
	keys.resetKeys()
	length := len(keys.keys)
	if length != 0 {
		t.Errorf("Expected length of keys to be 0, instead found %d", length)
	}
}

func TestAddKey(t *testing.T) {
	var keys hotKeys
	buildKeys(&keys)
	key := hotKey{key: tcell.KeyEnd, action: func() {}, text: "text4"}
	keys.addKey(key)
	length := len(keys.keys)
	if length != 5 {
		t.Errorf("Expected length of keys to be 5, instead found %d", length)
	}
}

func TestGenerateHelp(t *testing.T) {
	var keys hotKeys
	buildKeys(&keys)
	help := keys.generateHelp()
	if help != "text0 | text1 | text2 | text3" {
		t.Errorf("Function is not generating help according to plan")
	}
}
