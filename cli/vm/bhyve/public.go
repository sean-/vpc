package bhyve

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type BhyveConfig struct {
	Name           string
	Uuid           string
	Vcpus          int
	Ram            string
	Diskdriver     string
	Diskdevice     string
	Disksize       string
	Nicdriver      string
	Nicdevice      string
	Serialconsole1 string
	Serialconsole2 string
}

func PrintConfig(config BhyveConfig) {
	log.Info().
		Str("name", config.Name).
		Str("uuid", config.Uuid).
		Int("vcpus", config.Vcpus).
		Str("ram", config.Ram).
		Str("diskdriver", config.Diskdriver).
		Str("diskdevice", config.Diskdevice).
		Str("disksize", config.Disksize).
		Str("nicdriver", config.Nicdriver).
		Str("nicdevice", config.Nicdevice).
		Str("serialconsole1", config.Serialconsole1).
		Str("serialconsole2", config.Serialconsole2).
		Msg("Creating Bhyve VM")

	fmt.Printf("name:          \t%s\n", config.Name)
	fmt.Printf("uuid:          \t%s\n", config.Uuid)
	fmt.Printf("vcpus:         \t%d\n", config.Vcpus)
	fmt.Printf("ram:           \t%s\n", config.Ram)
	fmt.Printf("diskdriver:    \t%s\n", config.Diskdriver)
	fmt.Printf("diskdevice:    \t%s\n", config.Diskdevice)
	fmt.Printf("disksize:      \t%s\n", config.Disksize)
	fmt.Printf("nicdriver:     \t%s\n", config.Nicdriver)
	fmt.Printf("nicdevice:     \t%s\n", config.Nicdevice)
	fmt.Printf("serialconsole1:\t%s\n", config.Serialconsole1)
	fmt.Printf("serialconsole2:\t%s\n", config.Serialconsole2)
}
