// Cloud phone call: rings a device with a spoken prompt. Uses the call key.
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
	// CallKey 留空時，會經 VIBLE_CALL_KEY 或本機 /api/calls/key 自動取得。
	res, err := c.Call("Remind me to drink water.",
		vible.WithTitle("Hydration reminder"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("call_id=%s ringing=%v delivered=%v\n", res.CallID, res.Ringing, res.Delivered)
}
