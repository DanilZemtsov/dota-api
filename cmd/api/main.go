package main

import (
	"log"
	"www/pkg/processor"
	"www/pkg/server"
	"www/pkg/storage"
)

func main() {

	db, err := storage.Connect("root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		log.Fatalf("Db error: %v", err)
	}
	defer db.Close()
	go processor.MegaUpdateHero(db)
	server.HandleRequests(db)

}
