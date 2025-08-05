package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Zigelzi/go-rss-gator/internal/config"
	"github.com/Zigelzi/go-rss-gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	currentConfig *config.Config
	db            *database.Queries
}

func main() {
	// Loading the configuration from file system.
	newConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", newConfig.DbURL)
	if err != nil {
		log.Fatalf("unable to open database connection: %w", err)
	}
	dbQueries := database.New(db)

	appState := state{
		currentConfig: &newConfig,
		db:            dbQueries,
	}

	// Registering the existing commands.
	cmds := commands{}

	// Debug
	cmds.register("reset", handleReset)

	// Users
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handleListUsers)

	// RSS feeds
	cmds.register("agg", handleAggregate)
	cmds.register("feeds", handleListFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("follow", middlewareLoggedIn(handleFollowFeed))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollowFeed))
	cmds.register("following", middlewareLoggedIn(handleListFollowedFeeds))
	cmds.register("browse", middlewareLoggedIn(handleBrowsePosts))

	if len(os.Args) < 2 {
		log.Fatal("command arguments are missing")
	}
	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(&appState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
