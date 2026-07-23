// Voice converse: send audio, the agent speaks its reply on the device.
package main

import (
	"fmt"
	"log"
	"os"

	vible "github.com/Algohash-inc/Vible/sdk/go"
)

func main() {
	path := "out.wav" // 用 speech 範例產生嘅檔做輸入；可用第一個參數覆蓋
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	audio, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v (run the speech example first)", path, err)
	}
	c, err := vible.New()
	if err != nil {
		log.Fatal(err)
	}
	res, err := c.Converse(audio)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("transcript: %s\nspoken: %v\n", res.Transcript, res.Spoken)
}
