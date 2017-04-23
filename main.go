package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/duskhacker/cqrsnu/api"
	"github.com/duskhacker/cqrsnu/cafe"
	"github.com/duskhacker/cqrsnu/internal/github.com/bitly/nsq/internal/app"
)

func main() {
	cafe.SetLookupdHTTPAddrs(app.StringArray{})
	cafe.Init()
	chef_todos.Init()
	api.GinEngine().Run(":8080")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			log.Println("Exiting")
			cafe.StopAllConsumers()
			chef_todos.StopAllConsumers()
			log.Println("Done")
			os.Exit(0)
		}

	}()
}
