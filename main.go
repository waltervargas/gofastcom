package main

import (
	"fmt"
	"log"

	// "github.com/waltervargas/gofastcom/fast"
	"github.com/waltervargas/gofastcom/fastapi"
)

// func main() {https://api.fast.com/netflix/speedtest/v2
// 	bps, err := fast.Measure()
// 	fmt.Println(bps/125000, "mbps", err)
// }

const (
	maxRequest = 1
	maxTime    = 30
)

func main() {
	client, err := fastclient.New(maxRequest, maxTime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)
}
