package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Block struct {
	FullText string `json:"full_text"`
}

func main() {

	// init
	out := bufio.NewWriter(os.Stdout)
	enc := json.NewEncoder(out)
	fmt.Fprintln(out, `{ "version": 1 }`)
	fmt.Fprintln(out, `[`)
	out.Flush()

	clock := &Block{}
	blocks := []*Block{
		clock,
	}

	for range time.Tick(time.Second) {

		clock.FullText = time.Now().Format("Mon 2 Jan 15:04:05 2006")

		enc.Encode(blocks)
		fmt.Fprintln(out, ",")
		out.Flush()
	}
}
