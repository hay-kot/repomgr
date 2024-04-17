package commander

type AppAction string

func (c AppAction) String() string {
	return string(c)
}

const (
	AppActionFork AppAction = ":GitFork"
	AppActionExit AppAction = ":Exit"
)
