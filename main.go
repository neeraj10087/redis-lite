package main

import (
	"fmt"
	"redis-lite/server"
	"redis-lite/store"
)

func main() {
	s := store.New()
	fmt.Println("starting redis-lite...")
	server.StartAsync(":7379", s)
}