package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"wie.gg/lnk/handler"
	"wie.gg/lnk/store"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".env not loaded: %v", err)
	}

	store.InitializeStore(nil)
}

func main() {

	r := handler.SetupRouter()

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
