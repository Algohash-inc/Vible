// Non-streaming chat completion.
package main

import (
	"fmt"
	"log"

	vible "github.com/Algohash-inc/Vible/sdk/go"
)

func main() {
	c, err := vible.New()
	if err != nil {
		log.Fatal(err)
	}
	reply, err := c.Chat([]vible.Message{
		{Role: "user", Content: "In one sentence, what is Vible?"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
