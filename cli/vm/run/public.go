package create

import (
	"fmt"
	"io/ioutil"

	zfs "github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	createKeyName           = "vmname"
	createKeyUUID           = "uuid"
	createKeyVCPUs          = "vcpus"
	createKeyRAM            = "ram"
	createKeyDiskDriver     = "diskdriver"
	createKeyDiskDevice     = "diskdevice"
	createKeyDiskSize       = "disksize"
	createKeyNicDriver      = "nicdriver"
	createKeyNicDevice      = "nicdevice"
	createKeySerialConsole1 = "serialconsole1"
	createKeySerialConsole2 = "serialconsole2"
)

type BhyveConfig struct {
	name           string
	uuid           string
	vcpus          int
	ram            string
	diskdriver     string
	diskdevice     string
	disksize       string
	nicdriver      string
	nicdevice      string
	serialconsole1 string
	serialconsole2 string
}

func printConfig(config BhyveConfig) {
	log.Info().
		Str("name", config.name).
		Str("uuid", config.uuid).
		Int("vcpus", config.vcpus).
		Str("ram", config.ram).
		Str("diskdriver", config.diskdriver).
		Str("diskdevice", config.diskdevice).
		Str("disksize", config.disksize).
		Str("nicdriver", config.nicdriver).
		Str("nicdevice", config.nicdevice).
		Str("serialconsole1", config.serialconsole1).
		Str("serialconsole2", config.serialconsole2).
		Msg("Creating Bhyve VM")

	fmt.Printf("name:          \t%s\n", config.name)
	fmt.Printf("uuid:          \t%s\n", config.uuid)
	fmt.Printf("vcpus:         \t%d\n", config.vcpus)
	fmt.Printf("ram:           \t%s\n", config.ram)
	fmt.Printf("diskdriver:    \t%s\n", config.diskdriver)
	fmt.Printf("diskdevice:    \t%s\n", config.diskdevice)
	fmt.Printf("disksize:      \t%s\n", config.disksize)
	fmt.Printf("nicdriver:     \t%s\n", config.nicdriver)
	fmt.Printf("nicdevice:     \t%s\n", config.nicdevice)
	fmt.Printf("serialconsole1:\t%s\n", config.serialconsole1)
	fmt.Printf("serialconsole2:\t%s\n", config.serialconsole2)
}

func getZpoolName() (string, error) {
	zpools, err := zfs.ListZpools()
	if err != nil {
		return nil, errors.Wrap(err, "unable to find zpool")
	}

	// Always returns the first zpool found
	zpoolName := zpools[0].Name
	log.Info().Str("zpool", zpoolName).Msg("Using zpool")

	return zpoolName, nil
}

func getGuestPath() string {
	zpoolName, err := getZpoolName()
	if err != nil {
		return errors.Wrap(err, "unable to get zpool name")
	}

	path := fmt.Sprintf("%s/%s", zpoolName, viper.GetString(createKeyUUID))

	return path
}

func getDiskPath() string {
	path := fmt.Sprintf("%s/disk0", getGuestPath())
	return path
}

// Create the datasets for bhyve
// Example from chyves
// /chyves/
//	`-- zones
//     |-- Firmware
//     |-- ISO
//     |   |-- null.iso
//     |   `-- ubuntu-16.04.3-server-amd64.iso
//	   |-- guests
//	   |   `-- transcode
//     |       |-- img
//     |       `-- logs
//     	`-- logs
func setupDatasets() error {
	fsName := getGuestPath()
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().
			Str("filesystem", fsName).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsName, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName = fmt.Sprintf("%s/firmware", getGuestPath())
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().
			Str("filesystem", fsName).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsName, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName = fmt.Sprintf("%s/iso", getGuestPath())
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().
			Str("filesystem", fsName).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsName, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName = fmt.Sprintf("%s/disk0", getGuestPath())
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().Str("volume", fsName).
			Str("disksize", viper.GetString(createKeyDiskSize)).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateVolume(fsName, viper.GetString(createKeyDiskSize), nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	return nil
}

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "create",
		Short:        "Create a Virtual Machine",
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "create").Msg("")

			cfg, err := config.New()
			if err != nil {
				return errors.Wrap(err, "unable to load configuration")
			}

			if err := cfg.Load(); err != nil {
				return errors.Wrapf(err, "unable to load %s config", buildtime.PROGNAME)
			}

			a, err := agent.New(cfg)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}
			defer a.Shutdown()

			// verify db credentials
			if err := a.Pool().Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}

			// Wrap jackc/pgx in an sql.DB-compatible facade.
			db, err := a.Pool().STDDB()
			if err != nil {
				return errors.Wrap(err, "unable to conjur up sql.DB facade")
			}

			if err := db.Ping(); err != nil {
				return errors.Wrap(err, "unable to ping with stdlib driver")
			}

			config := BhyveConfig{
				name:           viper.GetString(createKeyName),
				uuid:           viper.GetString(createKeyUUID),
				vcpus:          viper.GetInt(createKeyVCPUs),
				ram:            viper.GetString(createKeyRAM),
				diskdriver:     viper.GetString(createKeyDiskDriver),
				diskdevice:     viper.GetString(createKeyDiskDevice),
				disksize:       viper.GetString(createKeyDiskSize),
				nicdriver:      viper.GetString(createKeyNicDriver),
				nicdevice:      viper.GetString(createKeyNicDevice),
				serialconsole1: viper.GetString(createKeySerialConsole1),
				serialconsole2: viper.GetString(createKeySerialConsole2),
			}
			printConfig(config)

			// Setup ZFS datasets and zvol
			if err := setupDatasets(); err != nil {
				return errors.Wrap(err, "Failed to setup ZFS datasets for virtual machine")
			}

			// Setup Networking
			// (if needed)

			// Write out device.map
			deviceMap := fmt.Sprintf("(hd0) %s", getDiskPath())
			deviceMapPath := fmt.Sprintf("%s/device.map", getGuestPath())
			if err := ioutil.WriteFile(deviceMapPath, deviceMap, 0644); err != nil {
				return errors.Wrap(err, "Cannot write device.map")
			}

			// Write out config to database

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		return nil
	},
}

func init() {

	{
		const (
			key          = createKeyName
			longName     = "name"
			shortName    = "n"
			defaultValue = ""
			description  = "Virtual Machine Name"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyUUID
			longName     = "uuid"
			shortName    = "u"
			defaultValue = ""
			description  = "UUID"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyVCPUs
			longName     = "vcpus"
			shortName    = "c"
			defaultValue = 1
			description  = "Number of vCPUs"
		)

		flags := Cmd.Cobra.Flags()
		flags.UintP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyRAM
			longName     = "ram"
			shortName    = "r"
			defaultValue = "256M"
			description  = "RAM"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyDiskDriver
			longName     = "diskdriver"
			shortName    = ""
			defaultValue = "virtio-blk"
			description  = "Disk Driver Emulation"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyDiskDevice
			longName     = "diskdevice"
			shortName    = "D"
			defaultValue = ""
			description  = "Path to disk image/block device"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyDiskSize
			longName     = "disksize"
			shortName    = ""
			defaultValue = "256M"
			description  = "Disk size in Megabytes"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyNicDriver
			longName     = "nicdriver"
			shortName    = ""
			defaultValue = "virtio-net"
			description  = "NIC Driver Emulation"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeyNicDevice
			longName     = "nicdevice"
			shortName    = "N"
			defaultValue = ""
			description  = "NIC Driver Emulation"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeySerialConsole1
			longName     = "serialconsole1"
			shortName    = "s"
			defaultValue = ""
			description  = "Serial Console 1 Device"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = createKeySerialConsole2
			longName     = "serialconsole2"
			shortName    = "a"
			defaultValue = ""
			description  = "Serial Console 2 Device"
		)

		flags := Cmd.Cobra.Flags()
		flags.StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, flags.Lookup(longName))
		viper.SetDefault(key, defaultValue)
	}
}
