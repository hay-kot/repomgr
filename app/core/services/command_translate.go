package services

import tea "github.com/charmbracelet/bubbletea"

func TranslateTeaKey(key tea.KeyType) (bool, string) {
	var keystr string
	switch key {
	case tea.KeyEnter:
		keystr = "enter"
	case tea.KeyCtrlO:
		keystr = "ctrl-o"
	case tea.KeyCtrlA:
		keystr = "ctrl-a"
	case tea.KeyCtrlB:
		keystr = "ctrl-b"
	case tea.KeyCtrlD:
		keystr = "ctrl-d"
	case tea.KeyCtrlE:
		keystr = "ctrl-e"
	case tea.KeyCtrlF:
		keystr = "ctrl-f"
	case tea.KeyCtrlG:
		keystr = "ctrl-g"
	case tea.KeyCtrlH:
		keystr = "ctrl-h"
	case tea.KeyCtrlI:
		keystr = "ctrl-i"
	case tea.KeyCtrlJ:
		keystr = "ctrl-j"
	case tea.KeyCtrlK:
		keystr = "ctrl-k"
	case tea.KeyCtrlL:
		keystr = "ctrl-l"
	case tea.KeyCtrlN:
		keystr = "ctrl-n"
	case tea.KeyCtrlP:
		keystr = "ctrl-p"
	case tea.KeyCtrlQ:
		keystr = "ctrl-q"
	case tea.KeyCtrlR:
		keystr = "ctrl-r"
	case tea.KeyCtrlS:
		keystr = "ctrl-s"
	case tea.KeyCtrlT:
		keystr = "ctrl-t"
	case tea.KeyCtrlU:
		keystr = "ctrl-u"
	case tea.KeyCtrlV:
		keystr = "ctrl-v"
	case tea.KeyCtrlW:
		keystr = "ctrl-w"
	case tea.KeyCtrlX:
		keystr = "ctrl-x"
	case tea.KeyCtrlY:
		keystr = "ctrl-y"
	case tea.KeyCtrlZ:
		keystr = "ctrl-z"
	default:
		return false, "" // Handle unknown keys
	}
	return true, keystr
}
