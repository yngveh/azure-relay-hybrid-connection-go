package main

import (
	"context"
	"log"
	"os"
	"os/signal"

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

	h := func(resp *arhc.Response, req *arhc.Request) error {
		log.Print("REQ", *req)
		resp.SetResponseBody([]byte("HELLO"))
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.Listen(ctx, h); err != nil {
		panic(err)
	}

	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt)
	<-e
}
