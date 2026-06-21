// Text-to-speech: synthesizes WAV bytes and writes out.wav.
package main

import (
	"log"
	"os"

	vible "github.com/Algohash-inc/Vible/sdk/go"
)

func main() {
	c, err := vible.New()
	if err != nil {
		log.Fatal(err)
	}
	// voice 留空 -> 預設聲。亦可用 vible.WithVoice("clone:NAME") 或內建名如 "nova"。
	wav, err := c.Speech("Hello from the Vible Go SDK.")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("out.wav", wav, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote out.wav (%d bytes)", len(wav))
}
