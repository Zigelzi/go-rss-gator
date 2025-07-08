package main

import (
	"fmt"
	"log"

	"github.com/Zigelzi/go-rss-gator/internal/config"
)

func main() {
	fmt.Println("This is a RSS feed aggregator!")
	newConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	newConfig.SetUser("miika")
	updatedConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(updatedConfig)
}
