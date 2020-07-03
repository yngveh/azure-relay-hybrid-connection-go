package main

import (
	"log"
	"os"

	arhc "github.com/yngveh/azure-relay-hybrid-connection-go"
)

func main() {

	// Read config from environment variables
	c := &arhc.Config{
		Namespace:      os.Getenv("HC_NAMESPACE"),
		ConnectionName: os.Getenv("HC_CONNECTION_NAME"),
		KeyName:        os.Getenv("HC_KEY_NAME"),
		Key:            os.Getenv("HC_KEY"),
	}

	// Create a azure relay hybrid connection client
	client, err := arhc.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Create a handler function
	h := func(resp *arhc.Response, req *arhc.Request) error {

		log.Print("request", req.Target)
		resp.SetResponseBody([]byte("HELLO"))
		return nil
	}

	// Start blocking listener
	if err := client.Listen(h); err != nil {
		panic(err)
	}

}
