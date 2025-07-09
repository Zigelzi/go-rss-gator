package main

import "errors"

// Validates that the login command contains the username argument
// and sets the current user to the one provided in the argument.
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("username is required argument in login command")
	}
	err := s.currentConfig.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	return nil
}
