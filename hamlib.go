package main

import (
	"github.com/dh1tw/goHamlib"
	"log"
)

func handleHamlib() {
	rig := goHamlib.Rig{}
	goHamlib.SetDebugLevel(goHamlib.DebugNone)
	if err := rig.Init(373); err != nil {
		panic(err)
	}
	if err := rig.SetPort(goHamlib.Port{
		RigPortType: goHamlib.RigPortSerial,
		Portname:    "/dev/ttyUSB0",
		Baudrate:    9600,
		Databits:    8,
		Stopbits:    1,
		Parity:      goHamlib.ParityNone,
		Handshake:   goHamlib.HandshakeNone,
	}); err != nil {
		panic(err)
	}
	if err := rig.Open(); err != nil {
		panic(err)
	}

	info, err := rig.GetFreq(goHamlib.VFOCurrent)
	if err != nil {
		panic(err)
	}
	log.Printf("Frequency: %f", info)
}
