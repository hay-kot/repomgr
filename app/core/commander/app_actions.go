package commander

import "strings"

type AppAction string

const (
	AppActionFork AppAction = ":GitFork"
	AppActionExit AppAction = ":Exit"
)

func (c AppAction) String() string {
	return string(c)
}

func ParseAppAction(s string) (a AppAction, arg string, ok bool) {
	if len(s) == 0 {
		return "", arg, false
	}

	if s[0] != ':' {
		return "", arg, false
	}

	parts := strings.Split(s, " ")
	a, arg = AppAction(parts[0]), strings.Join(parts[1:], " ")

	if !a.IsValid() {
		return "", "", false
	}

	return a, arg, true
}

func (c AppAction) IsValid() bool {
	switch c {
	case AppActionFork, AppActionExit:
		return true
	default:
		return false
	}
}
