package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	red  = "#ff0000"
	gray = "#777777"
)

var (
	warnBat     = flag.Float64("warn-bat", 0, "low battery warning percentage")
	warnWatts   = flag.Float64("warn-power", 0, "high power draw warning (watts)")
	batteryPath = flag.String("battery", "/sys/class/power_supply/BAT0", "battery device")
)
var (
	capacity   = "/sys/class/power_supply/BAT0/capacity"
	currentNow = "/sys/class/power_supply/BAT0/current_now"
	voltageNow = "/sys/class/power_supply/BAT0/voltage_now"
	status     = "/sys/class/power_supply/BAT0/status"
)

type Block struct {
	FullText string `json:"full_text"`
	Color    string `json:"color,omitempty"`
}

func main() {

	flag.Parse()
	capacity = path.Join(*batteryPath, "capacity")
	currentNow = path.Join(*batteryPath, "current_now")
	voltageNow = path.Join(*batteryPath, "voltage_now")
	status = path.Join(*batteryPath, "status")

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
		discharging := batteryDischarging()
		power.FullText = fmt.Sprintf("% 6.2f W ", watts)
		const t = 0.95 // slow recursive filter
		if discharging {
			avgWatts = (t * avgWatts) + (1-t)*watts
		} else {
			avgWatts = 0
		}

		switch {
		case *warnWatts != 0 && avgWatts > *warnWatts && discharging:
			power.Color = red
		case !discharging:
			power.Color = gray
		default:
			power.Color = ""
		}

		// battery capacity
		pct := batteryPct()
		bat.FullText = fmt.Sprintf("% 4.0f %% ", pct)
		switch {
		case pct <= *warnBat:
			bat.Color = red
		case !discharging:
			bat.Color = gray
		default:
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
	microVolt := readFloat64(voltageNow)
	microAmp := readFloat64(currentNow)
	return microVolt * microAmp / 1e12
}

func batteryPct() float64 {
	return readFloat64(capacity)
}

func batteryDischarging() bool {
	return readString(status) == "Discharging"
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
