package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (cmds *commands) register(cmdName string, handler func(*state, command) error) {
	if cmds.registeredCommands == nil {
		cmds.registeredCommands = make(map[string]func(*state, command) error)
	}
	// Does this need to be made thread safe with Mutex?
	cmds.registeredCommands[cmdName] = handler

}

func (cmds *commands) run(s *state, cmd command) error {
	handler, exists := cmds.registeredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("command with name '%s' doesn't exist", cmd.Name)
	}
	return handler(s, cmd)
}
