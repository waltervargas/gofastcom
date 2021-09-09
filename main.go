package main

import (
	"fmt"
	"log"

	"github.com/waltervargas/gofastcom/fastapi"
)

const (
	maxRequest = 1
	maxTime    = 30
)

func main() {
	client, err := fastclient.New(maxRequest, maxTime)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range client.OCA.Targets {
		fmt.Printf("city: %s, url: %s\n", c.Location.City, c.URL)
	}
}
