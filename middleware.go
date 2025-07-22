package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zigelzi/go-rss-gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.currentConfig.CurrentUserName)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("you need to be logged in to user command '%s'", cmd.Name)
			}
			return fmt.Errorf("unable to get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
