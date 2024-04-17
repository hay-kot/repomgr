package ui

// state is a struct that holds some generic state information for the UI.
// Primarily, it's purpose is to maintain a consisten contract between the
// outer layers of the application and the UI componenents. This is useful
// for keeping patterns consistent between pages/components.
type state struct {
	err error
	msg string

	exit bool
}

func (s *state) ExitMessage() string {
	return s.msg
}

// isExit indicates that the UI should exit when possible
func (s *state) isExit() bool {
	return false
}

// recieveError sets the error state based on a channel when recieved.
// This is useful for when a command is running in the background and
// an error is recieved. This will spawn a goroutine to wait for the error
func (s *state) recieveError(errch <-chan error) {
	go func() {
		s.err = <-errch
	}()
}

func (s *state) signalExit(msg string) {
	s.msg = msg
	s.exit = true
}
