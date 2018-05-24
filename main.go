package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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
	power := &Block{}

	blocks := []*Block{
		power,
		clock,
	}

	for range time.Tick(time.Second) {

		clock.FullText = time.Now().Format("Mon 2 Jan 15:04:05 2006")

		power.FullText = fmt.Sprintf("%4.2f W ", batteryWatts())

		enc.Encode(blocks)
		fmt.Fprintln(out, ",")
		out.Flush()
	}
}

func batteryWatts() float64 {
	microVolt := readFloat64("/sys/class/power_supply/BAT0/voltage_now")
	microAmp := readFloat64("/sys/class/power_supply/BAT0/current_now")
	return microVolt * microAmp / 1e12
}

func readFloat64(file string) float64 {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return 0
	}
	str := strings.TrimSpace(string(bytes))
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println(err)
	}
	return v
}
