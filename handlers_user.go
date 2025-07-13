package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Zigelzi/go-rss-gator/internal/database"
	"github.com/google/uuid"
)

// Validates that the login command contains the username argument
// and sets the current user to the one provided in the argument.
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("username is required argument")
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("unable to query for existing user: %w", err)
	}
	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("user with username %s doesn't exist ", name)
	}

	err = s.currentConfig.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	return nil
}

// Registers new user to the service.
// Username is passed as command argument and needs to be unique.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("username is required argument")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("got over 1 argument (%d): %v", len(cmd.Args), cmd.Args)
	}

	name := cmd.Args[0]
	existingUser, err := s.db.GetUser(context.Background(), name)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("unable to query for existing user: %w", err)
	}
	emptyUser := database.User{}
	if existingUser != emptyUser {
		return fmt.Errorf("user with username %s already exists", name)
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	createdUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("unable to create user %w", err)
	}

	log.Println("created user:")
	log.Printf("%v", createdUser)
	err = s.currentConfig.SetUser(createdUser.Name)
	if err != nil {
		return fmt.Errorf("unable to set user to %s: %w", createdUser.Name, err)
	}
	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to delete users: %w", err)
	}
	return nil
}

// Lists all users that are registered to the service.
// Highlights the currently logged in user with (current)
func handleListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}
	if len(users) == 0 {
		fmt.Println("No registered users.")
		return nil
	}

	fmt.Println("Users:")
	for i, user := range users {
		printedText := "* " + user.Name
		if user.Name == s.currentConfig.CurrentUserName {
			printedText += " (current)"
		}
		if i != len(users) {
			printedText += "\n"
		}
		fmt.Print(printedText)
	}
	return nil
}
