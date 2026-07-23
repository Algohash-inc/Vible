// Streaming chat completion (SSE): prints tokens as they arrive.
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
	err = c.ChatStream([]vible.Message{
		{Role: "user", Content: "Count from 1 to 5, one number per line."},
	}, func(delta string) {
		fmt.Print(delta)
		os.Stdout.Sync()
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
