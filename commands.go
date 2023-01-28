package main

const (
	SUB   int = 1
	PUB   int = 2
	LIST  int = 3
	UNSUB int = 4
	QUIT  int = 5
)

// Command
// Id : unique id for the command
// Emptor : client which sent the command
// Args : arguments to the command (see docs.)
type Command struct {
	id     int
	emptor *Emptor
	args   []string
}
