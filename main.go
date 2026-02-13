package main

import (
	"flag"
	"log"

	config "github.com/anujkaushik1/GoDis/Config"
	"github.com/anujkaushik1/GoDis/sever"
)

func setupFlags() {

	flag.StringVar(&config.App.Host, "host", "0.0.0.0", "Host address for the server")
	flag.IntVar(&config.App.Port, "port", 7379, "Port number for the server")
	flag.Parse()

}

func main() {
	log.Println("Server is starting...")
	setupFlags()

	sever.RunTcpServer()
}
