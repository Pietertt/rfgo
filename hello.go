package main

import (
    "log"
    "github.com/brutella/hc"
    "github.com/brutella/hc/accessory"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
	// "math/rand"
	// "math"
)

var transmitter rpio.Pin
var led rpio.Pin

var mod float64 = 0.90

var superLong float64 = 2610 * mod; // 2610
var veryLong float64 = 1352 * mod; // 1352
var top float64 = 220 * mod; //220
var bottom float64 = 330 * mod; //330

var sendOn = [] int {
	0, 2, 2, 1, 2, 3, 1, 2, 2, 2,
	2, 2, 3, 1, 2, 2, 3, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 3,
	2, 1, 3, 1 }

var sendOff = [] int {
	0, 2, 2, 1, 2, 3, 1, 2, 2, 2,
	2, 2, 3, 1, 2, 2, 3, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 1, 3, 1 }

func main() {
	err := rpio.Open()

	led = rpio.Pin(4)
	transmitter = rpio.Pin(17)

	transmitter.Output()
	led.Output()

	if err != nil {
		panic(err)
	}

    info := accessory.Info{
		Name: "Lamp",
		SerialNumber: "",
		Manufacturer: "Pieter Boersma",
		Model: "1.0",
		FirmwareRevision: "1.0.1",
	}

    bedroomLamp := accessory.NewSwitch(info)

	bedroomLamp.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {	
			for i := 0; i < 20; i++ {
				sendRF(sendOn)
				time.Sleep(13 * time.Millisecond)
				led.High()
				time.Sleep(13 * time.Millisecond)
				led.Low()
			}
		} else {
			for i := 0; i < 20; i++ {
				sendRF(sendOff)
				time.Sleep(25 * time.Millisecond)
				led.High()
				time.Sleep(25 * time.Millisecond)
				led.Low()
			}
		}
	})

    configuration := hc.Config{Pin: "00102003"}
    t, err := hc.NewIPTransport(configuration, bedroomLamp.Accessory)
    if err != nil {
        log.Panic(err)
    }

    hc.OnTermination(func(){
        <-t.Stop()
    })

    t.Start()
}

func sendRF(values []int) {
	superLongs := time.Duration(superLong * 1000)
	veryLongs := time.Duration(veryLong * 1000)
	tops := time.Duration(top * 1000)
	bottoms := time.Duration(bottom * 1000)

	for i := 0; i < 34; i++ {
		if values[i] == 0 {
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(superLongs)
		}

		if values[i] == 1 {
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(veryLongs)
		}

		if values[i] == 2 {
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(bottoms)
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(veryLongs)
		}

		if values[i] == 3 {
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(bottoms)
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(bottoms)
			transmitter.High()
			time.Sleep(tops)
			transmitter.Low()
			time.Sleep(veryLongs)
		}
	}
}
