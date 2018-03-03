package list

import (
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName = "list"
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

// func listVMs() ([]string, error) {
// 	guestPath, err := bhyve.GetGuestPath()
// 	if err != nil {
// 		return errors.Wrap(err, "unable to get guest path")
// 	}

// 	files, err := ioutil.ReadDir(guestPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, f := range files {
// 		fmt.Println(f.Name())
// 	}

// 	return
// }

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

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "List Virtual Machines",
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()
			log.Info().Str("command", "list").Msg("")

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

			// TODO: Once DB work is done, we can query the db for this information
			// for now we will use the filesystem. We're defining inactive/"off" vms
			// to be vms that have a directory
			table := tablewriter.NewWriter(cons)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(false)
			table.SetAutoFormatHeaders(true)

			table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")

			table.SetHeader([]string{"name", "status"})

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
