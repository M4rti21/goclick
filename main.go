package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bendahl/uinput"
)

type Args struct {
	countdown int64
	duration  int64
	rate      int64
	button    string
}

func main() {
	var args Args
	flag.Int64Var(&args.countdown, "c", 3000, "Countdown (ms), default 3000")
	flag.Int64Var(&args.duration, "d", 10000, "Duration (ms), default 10000")
	flag.Int64Var(&args.rate, "r", 25, "Click Rate (ms), default 25")
	flag.StringVar(&args.button, "b", "left", "Button to emulate: left/l (default), right/r, middle/m")

	flag.Parse()

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("Virtual Mouse"))
	if err != nil {
		panic(err)
	}

	var click func() error

	switch args.button {
	case "l", "left":
		click = mouse.LeftClick
		break
	case "r", "right":
		click = mouse.RightClick
		break
	case "m", "middle":
		click = mouse.RightClick
		break
	default:
		panic(fmt.Sprintf("Invalid button: %s", args.button))
	}

	rate := time.Duration(args.rate) * time.Millisecond

	time.Sleep(time.Duration(args.countdown) * time.Millisecond)

	if args.duration <= 0 {
		for {
			err = click()
			if err != nil {
				panic(err)
			}
			time.Sleep(rate)
		}
	} else {
		time_left := time.Duration(args.duration) * time.Millisecond
		for {
			err = click()
			if err != nil {
				panic(err)
			}
			time_left -= rate
			if time_left <= 0 {
				break
			}
			time.Sleep(rate)
		}
	}

	err = mouse.Close()
	if err != nil {
		panic(err)
	}
}
