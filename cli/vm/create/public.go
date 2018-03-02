package create

import (
	"fmt"
	"io/ioutil"

	zfs "github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/cli/vm/bhyve"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName = "create"
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

func getZpoolName() (string, error) {
	zpools, err := zfs.ListZpools()
	if err != nil {
		return "", errors.Wrap(err, "unable to find zpool")
	}

	// Always returns the first zpool found
	zpoolName := zpools[0].Name
	log.Info().Str("zpool", zpoolName).Msg("Using zpool")

	return zpoolName, nil
}

func getGuestPath() (string, error) {
	zpoolName, err := getZpoolName()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", zpoolName, viper.GetString(createKeyUUID))

	return path, nil
}

func getDiskPath() (string, error) {
	guestPath, err := getGuestPath()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/disk0", guestPath)
	return path, nil
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
	fsRoot, err := getGuestPath()
	if err != nil {
		return err
	}

	if _, err := zfs.GetDataset(fsRoot); err != nil {
		log.Info().
			Str("filesystem", fsRoot).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsRoot, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName := fmt.Sprintf("%s/firmware", fsRoot)
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().
			Str("filesystem", fsName).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsName, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName = fmt.Sprintf("%s/iso", fsRoot)
	if _, err := zfs.GetDataset(fsName); err != nil {
		log.Info().
			Str("filesystem", fsName).
			Str("UUID", viper.GetString(createKeyUUID)).
			Msg("Creating ZFS Filesystem")
		if _, err := zfs.CreateFilesystem(fsName, nil); err != nil {
			return errors.Wrap(err, "unable to create vm filesystem")
		}
	}

	fsName = fmt.Sprintf("%s/disk0", fsRoot)
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
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
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

			config := bhyve.BhyveConfig{
				Name:           viper.GetString(createKeyName),
				Uuid:           viper.GetString(createKeyUUID),
				Vcpus:          viper.GetInt(createKeyVCPUs),
				Ram:            viper.GetString(createKeyRAM),
				Diskdriver:     viper.GetString(createKeyDiskDriver),
				Diskdevice:     viper.GetString(createKeyDiskDevice),
				Disksize:       viper.GetString(createKeyDiskSize),
				Nicdriver:      viper.GetString(createKeyNicDriver),
				Nicdevice:      viper.GetString(createKeyNicDevice),
				Serialconsole1: viper.GetString(createKeySerialConsole1),
				Serialconsole2: viper.GetString(createKeySerialConsole2),
			}
			bhyve.PrintConfig(config)

			// Setup ZFS datasets and zvol
			if err := setupDatasets(); err != nil {
				return errors.Wrap(err, "Failed to setup ZFS datasets for virtual machine")
			}

			// Setup Networking
			// (if needed)

			// Write out device.map
			diskPath, err := getDiskPath()
			if err != nil {
				return errors.Wrap(err, "unable to get vm path")
			}

			guestPath, err := getGuestPath()
			if err != nil {
				return errors.Wrap(err, "unable to get guest path")
			}

			deviceMap := fmt.Sprintf("(hd0) %s", diskPath)
			deviceMapPath := fmt.Sprintf("%s/device.map", guestPath)
			if err := ioutil.WriteFile(deviceMapPath, []byte(deviceMap), 0644); err != nil {
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
