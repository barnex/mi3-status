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

const (
	red = "#ff0000"
)

const (
	warnBat   = 30 // warn when battery below this percentage
	warnWatts = 8  // warn when over this power draw for an extended time
)

type Block struct {
	FullText string `json:"full_text"`
	Color    string `json:"color,omitempty"`
}

func main() {

	// init
	out := bufio.NewWriter(os.Stdout)
	enc := json.NewEncoder(out)
	fmt.Fprintln(out, `{ "version": 1 }`)
	fmt.Fprintln(out, `[`)
	out.Flush()

	clock := &Block{}
	bat := &Block{}
	power := &Block{}

	blocks := []*Block{
		&Block{}, // empty, separator
		bat,
		power,
		clock,
		&Block{}, // empty, separator
	}

	avgWatts := 0.0

	update := func(now time.Time) {
		// clock
		clock.FullText = now.Format(" Mon 2 Jan 15:04:05 2006 ")

		// power use
		watts := batteryWatts()
		power.FullText = fmt.Sprintf("% 6.2f W ", watts)
		const t = 0.95 // slow recursive filter
		avgWatts = (t * avgWatts) + (1-t)*watts
		if avgWatts > warnWatts && batteryDischarging() {
			power.Color = red
		} else {
			power.Color = ""
		}

		// battery capacity
		pct := batteryPct()
		bat.FullText = fmt.Sprintf("% 5.1f %% ", pct)
		if pct < warnBat {
			bat.Color = red
		} else {
			bat.Color = ""
		}

		// output
		enc.Encode(blocks)
		fmt.Fprintln(out, ",")
		out.Flush()
	}

	update(time.Now())
	for now := range time.Tick(time.Second) {
		update(now)
	}
}

func batteryWatts() float64 {
	microVolt := readFloat64("/sys/class/power_supply/BAT0/voltage_now")
	microAmp := readFloat64("/sys/class/power_supply/BAT0/current_now")
	return microVolt * microAmp / 1e12
}

func batteryPct() float64 {
	now := readFloat64("/sys/class/power_supply/BAT0/charge_now")
	full := readFloat64("/sys/class/power_supply/BAT0/charge_full")
	return 100 * now / full
}

func batteryDischarging() bool {
	return readString("/sys/class/power_supply/BAT0/voltage_now") == "Discharging"
}

func readFloat64(file string) float64 {
	str := readString(file)
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println(err)
	}
	return v
}

func readString(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return ""
	}
	return strings.TrimSpace(string(bytes))
}
