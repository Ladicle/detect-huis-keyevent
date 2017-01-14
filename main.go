package main

import (
	"fmt"
	"os"

	"github.com/gvalkov/golang-evdev"
	"net/http"
	"net/url"
)

const eventURL = "http://10.0.1.100:1880/huis"

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("You must specify device name.")
		os.Exit(1)
	}

	dev, err := evdev.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to open input device: %s\n", os.Args[1])
		os.Exit(1)
	}

	fmt.Printf("Listening for evens of %s.\n", os.Args[1])
	for {
		events, err := dev.Read()
		if err != nil {
			fmt.Printf("Can not read device envets: %v\n", err)
			continue
		}
		for i := range events {
			str, err := detectDownKeyEvent(&events[i])
			if err != nil {
				continue
			}
			fmt.Println(str)
			if err := sendKey(str); err != nil {
				fmt.Printf("Can not send key: %v\n", err)
			}
		}
	}
}

func sendKey(key string) error {
	values := url.Values{}
	values.Add("key", key)
	_, err := http.PostForm(eventURL, values)
	return err
}

func detectDownKeyEvent(ev *evdev.InputEvent) (string, error) {
	if ev.Value != 1 {
		return "", fmt.Errorf("Value is not 1: %d", ev.Value)
	}

	var codeName string
	code := int(ev.Code)

	val, hashkey := evdev.KEY[code]
	if hashkey {
		codeName = val
	} else {
		val, hashkey := evdev.BTN[code]
		if hashkey {
			codeName = val
		} else {
			codeName = "?"
		}
	}
	return codeName, nil
}
