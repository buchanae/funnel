// Copied and modified from: https://github.com/bcicen/ctop
// MIT License - Copyright (c) 2017 VektorLab

package termdash

import (
	ui "github.com/gizak/termui"
)

// Common action keybindings
var keyMap = map[string][]string{
	"up": {
		"/sys/kbd/<up>",
		"/sys/kbd/k",
	},
	"down": {
		"/sys/kbd/<down>",
		"/sys/kbd/j",
	},
	"left": {
		"/sys/kbd/<left>",
		"sys/kbd/h",
	},
	"right": {
		"/sys/kbd/<right>",
		"sys/kbd/l",
	},
	"exit": {
		"/sys/kbd/q",
		"/sys/kbd/C-c",
		"/sys/kbd/<escape>",
	},
	"help": {
		"/sys/kbd/h",
		"/sys/kbd/?",
	},
}

// Apply a common handler function to all given keys
func HandleKeys(i string, f func()) {
	for _, k := range keyMap[i] {
		ui.Handle(k, func(ui.Event) { f() })
	}
}
