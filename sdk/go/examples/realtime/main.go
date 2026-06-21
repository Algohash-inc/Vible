// Realtime over WebSocket: streams tokens via the hand-rolled RFC 6455 client.
package main

import (
	"fmt"
	"log"
	"os"

	vible "github.com/Algohash-inc/Vible/sdk/go"
)

func main() {
	c, err := vible.New()
	if err != nil {
		log.Fatal(err)
	}
	err = c.Realtime([]vible.Message{
		{Role: "user", Content: "Say hi in five words."},
	}, func(delta string) {
		fmt.Print(delta)
		os.Stdout.Sync()
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
