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

// MatchesActionSyntax checks if a string matches the syntax for an app action.
// this is a colon followed by a command.
func MatchesActionSyntax(s string) bool {
	return len(s) > 0 && s[0] == ':'
}

func ParseAppAction(s string) (a AppAction, rest string, ok bool) {
	if len(s) == 0 {
		return "", "", false
	}

	if s[0] != ':' {
		return "", "", false
	}

	parts := strings.Split(s, " ")
	a, rest = AppAction(parts[0]), strings.Join(parts[1:], " ")

	if !a.IsValid() {
		return "", "", false
	}

	return a, rest, true
}

func (c AppAction) IsValid() bool {
	switch c {
	case AppActionFork, AppActionExit:
		return true
	default:
		return false
	}
}
