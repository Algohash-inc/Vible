// Speech-to-text: transcribes out.wav (or the path in os.Args[1]).
package main

import (
	"fmt"
	"log"
	"os"

	vible "github.com/Algohash-inc/Vible/sdk/go"
)

func main() {
	path := "out.wav" // 預設讀 speech 範例產生嘅檔；可用第一個參數覆蓋
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
	text, err := c.Transcribe(audio)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(text)
}
