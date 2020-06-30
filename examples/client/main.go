package main

import (
	"log"
	"os"

	arhc "github.com/yngveh/azure-relay-hybrid-connection-go"
)

func main() {

	log.Print("START")
	defer log.Print("END")

	c := &arhc.Config{
		Namespace:      os.Getenv("HC_NAMESPACE"),
		ConnectionName: os.Getenv("HC_CONNECTION_NAME"),
		KeyName:        os.Getenv("HC_KEY_NAME"),
		Key:            os.Getenv("HC_KEY"),
	}

	client, err := arhc.NewClient(c)
	if err != nil {
		panic(err)
	}

	r, err := client.Get("/test")
	if err != nil {
		panic(err)
	}

	log.Print("RESPONSE ==> ", string(r))
}
