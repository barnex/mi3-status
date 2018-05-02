package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Block struct {
	FullText string `json:"full_text"`
}

func main() {
	fmt.Println(`{ "version": 1 }`)
	fmt.Println(`[`)

	clock := &Block{}

	blocks := []*Block{
		clock,
	}

	enc := json.NewEncoder(os.Stdout)
	for range time.Tick(time.Second) {
		clock.FullText = time.Now().Format("Mon 2 Jan 15:04:05 2006")
		enc.Encode(blocks)
		fmt.Println(",")
	}
}
