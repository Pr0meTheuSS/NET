package main

import (
	"log"
	"os"
	"snake/server"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	servePort, err := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 32)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := server.NewServer(int(servePort))
	log.Fatal(server.ListenAndServe())
}
