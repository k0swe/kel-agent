//go:build !hamlib

package config

func defaultHamlibConf() HamlibConfig {
	return HamlibConfig{
		Enabled:      false,
		RetrySeconds: 10,
		RigModel:     3073, // Icom IC-7300
		RigPort:      "RIG_PORT_SERIAL",
		PortName:     "/dev/ttyUSB0",
		BaudRate:     9600,
		DataBits:     8,
		StopBits:     1,
		Parity:       0, // None
		Handshake:    0, // None
	}
}
