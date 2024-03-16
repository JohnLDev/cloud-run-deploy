package main

import (
	"github.com/johnldev/4-deploy-cloud-run/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server.StartServer()
}
