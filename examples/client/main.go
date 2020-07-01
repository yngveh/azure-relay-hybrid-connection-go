package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

	// Use client.NewRequest instead of http.NewRequest
	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		panic(err)
	}

	// Rest is treated as normal http calls
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Print("RESPONSE ==> ", string(body))
}
