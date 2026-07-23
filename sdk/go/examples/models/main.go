// List available models.
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
	ids, err := c.Models()
	if err != nil {
		log.Fatal(err)
	}
	for _, id := range ids {
		fmt.Println(id)
	}
}
