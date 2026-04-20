//go:build hamlib

package config

import "github.com/xylo04/goHamlib"

func defaultHamlibConf() HamlibConfig {
	return HamlibConfig{
		Enabled:      false,
		RetrySeconds: 10,
		RigModel:     3073, // Icom IC-7300
		RigPort:      goHamlib.RigPortName[goHamlib.RigPortSerial],
		PortName:     "/dev/ttyUSB0",
		BaudRate:     9600,
		DataBits:     8,
		StopBits:     1,
		Parity:       byte(goHamlib.ParityNone),
		Handshake:    byte(goHamlib.HandshakeNone),
	}
}
