package main

import (
	"log"
	"os"

	"github.com/Zigelzi/go-rss-gator/internal/config"
)

type state struct {
	currentConfig *config.Config
}

func main() {
	// Loading the configuration from file system.
	newConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	appState := state{
		currentConfig: &newConfig,
	}

	// Registering the existing commands.
	cmds := commands{}
	cmds.register("login", handlerLogin)

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
