package main

import (
	"fmt"

	"wie.gg/lnk/handler"
	"wie.gg/lnk/store"
)

var env = "./.env"

func init() {
	store.InitializeStore(nil)
}

func main() {

	r := handler.SetupRouter(&env)

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
