package start

import (
	zfs "github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cli/vm/bhyve"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName                              = "start"
	_KeyISO                               = config.KeyISO
	_KeyJSON                              = config.KeyJson
	_KeyName                              = config.KeyName
	_KeyUUID                              = config.KeyUUID
	_KeyVCPUs                             = config.KeyVCPUs
	_KeyBootPartition                     = config.KeyBootPartition
	_KeyRAM                               = config.KeyRAM
	_KeyDiskDriver                        = config.KeyDiskDriver
	_KeyDiskDevice                        = config.KeyDiskDevice
	_KeyDiskSize                          = config.KeyDiskSize
	_KeyNicDriver                         = config.KeyNicDriver
	_KeyNicDevice                         = config.KeyNicDevice
	_KeyNicID                             = config.KeyNicID
	_KeySerialConsole1                    = config.KeySerialConsole1
	_KeySerialConsole2                    = config.KeySerialConsole2
	_KeyHostBridge                        = config.KeyHostBridge
	_KeyLPC                               = config.KeyLPC
	_KeyBhyveGenACPITables                = config.KeyBhyveGenACPITables
	_KeyBhyveIncGuestCoreMem              = config.KeyBhyveIncGuestCoreMem
	_KeyBhyveExitOnUnemuIOPort            = config.KeyBhyveExitOnUnemuIOPort
	_KeyBhyveYieldCPUOnHLT                = config.KeyBhyveYieldCPUOnHLT
	_KeyBhyveIgnoreUnimplementedMSRAccess = config.KeyBhyveIgnoreUnimplementedMSRAccess
	_KeyBhyveForceMSIInterrupts           = config.KeyBhyveForceMSIInterrupts
	_KeyBhyveAPICx2Mode                   = config.KeyBhyveAPICx2Mode
	_KeyBhyveDisableMPTableGeneration     = config.KeyBhyveDisableMPTableGeneration
	_KeyBhyveExitOnPause                  = config.KeyBhyveExitOnPause
	_KeyBhyveWireGuestMemory              = config.KeyBhyveWireGuestMemory
)

func checkVMExists(uuid string) error {
	fsName, err := bhyve.GetGuestPath(uuid)
	if err != nil {
		return errors.Wrap(err, "Unable to get guest dataset")
	}

	if _, err := zfs.GetDataset(fsName); err != nil {
		return errors.Wrap(err, "Unable to find guest")
	}

	return nil
}

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "Start a Virtual Machine",
		Aliases:      []string{"run"},
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString(_KeyUUID) == "" {
				return errors.Errorf("Must specify VM by UUID")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "run").Msg("")

			// cfg, err := config.New()
			// if err != nil {
			// 	return errors.Wrap(err, "unable to load configuration")
			// }

			// if err := cfg.Load(); err != nil {
			// 	return errors.Wrapf(err, "unable to load %s config", buildtime.PROGNAME)
			// }

			// a, err := agent.New(cfg)
			// if err != nil {
			// 	return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			// }
			// defer a.Shutdown()

			// // verify db credentials
			// if err := a.Pool().Ping(); err != nil {
			// 	return errors.Wrap(err, "unable to ping database")
			// }

			// // Wrap jackc/pgx in an sql.DB-compatible facade.
			// db, err := a.Pool().STDDB()
			// if err != nil {
			// 	return errors.Wrap(err, "unable to conjur up sql.DB facade")
			// }

			// if err := db.Ping(); err != nil {
			// 	return errors.Wrap(err, "unable to ping with stdlib driver")
			// }
			// Check if VM exists
			if err := checkVMExists(viper.GetString(config.KeyUUID)); err != nil {
				return errors.Wrap(err, "A VM with that UUID does not exist")
			}

			// Read in JSON config
			// Check if VM exists
			cfg, err := bhyve.ReadConfig(viper.GetString(config.KeyUUID))
			if err != nil {
				return errors.Wrap(err, "Cannot Read Config")
			}

			// cfg := bhyve.BhyveConfig{
			// 	Name:                         viper.GetString(_KeyName),
			// 	UUID:                         viper.GetString(_KeyUUID),
			// 	VCPUs:                        viper.GetInt(_KeyVCPUs),
			//  BootPartition:                          viper.GetString(_KeyBootPartition),
			// 	RAM:                          viper.GetString(_KeyRAM),
			// 	DiskDriver:                   viper.GetString(_KeyDiskDriver),
			// 	DiskDevice:                   viper.GetString(_KeyDiskDevice),
			// 	DiskSize:                     viper.GetString(_KeyDiskSize),
			// 	NicDriver:                    viper.GetString(_KeyNicDriver),
			// 	NicDevice:                    viper.GetString(_KeyNicDevice),
			// 	SerialConsole1:               viper.GetString(_KeySerialConsole1),
			// 	SerialConsole2:               viper.GetString(_KeySerialConsole2),
			// 	GenACPITables:                viper.GetBool(_KeyBhyveGenACPITables),
			// 	IncGuestCoreMem:              viper.GetBool(_KeyBhyveIncGuestCoreMem),
			// 	ExitOnUnemuIOPort:            viper.GetBool(_KeyBhyveExitOnUnemuIOPort),
			// 	YieldCPUOnHLT:                viper.GetBool(_KeyBhyveYieldCPUOnHLT),
			// 	IgnoreUnimplementedMSRAccess: viper.GetBool(_KeyBhyveIgnoreUnimplementedMSRAccess),
			// 	ForceMSIInterrupts:           viper.GetBool(_KeyBhyveForceMSIInterrupts),
			// 	Apicx2Mode:                   viper.GetBool(_KeyBhyveAPICx2Mode),
			// 	DisableMPTableGeneration:     viper.GetBool(_KeyBhyveDisableMPTableGeneration),
			// }
			bhyve.PrintConfig(*cfg)

			// Setup ZFS datasets and zvol
			// if err := setupDatasets(); err != nil {
			// 	return errors.Wrap(err, "Failed to setup ZFS datasets for virtual machine")
			// }

			// Setup Networking
			// (if needed)

			// Run grub-bhyve
			grubBhyve, err := bhyve.BuildGrubBhyveArgs(*cfg)
			if err != nil {
				return errors.Wrap(err, "unable to create bhyve command")
			}
			err = bhyve.RunGrubBhyve(*cfg, grubBhyve)
			if err != nil {
				return errors.Wrap(err, "unable to run grub-bhyve")
			}
			// Finally, run Bhyve
			bhyveString, err := bhyve.BuildBhyveArgs(*cfg)
			if err != nil {
				return errors.Wrap(err, "unable to create bhyve command")
			}
			err = bhyve.RunBhyve(*cfg, bhyveString)
			if err != nil {
				return errors.Wrap(err, "unable to run bhyve")
			}

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddStringFlag(self, _KeyISO, "iso", "", "", "ISO to be booted.", false, true); err != nil {
			return errors.Wrap(err, "unable to register iso flag on VPC VM create")
		}
		if err := flag.AddStringFlag(self, _KeyHostBridge, "hostbridge", "", "hostbridge", "A simple hostbridge. The amd_hostbridge emulation is identical but uses a PCI vendor ID of AMD.", false, true); err != nil {
			return errors.Wrap(err, "unable to register hostbridge flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyLPC, "lpc", "", "lpc", "Allow devices behind the LPC PCI-ISA bridge to be configured. The only supported devices are the TTY-class devices com1 and com2 and the boot ROM device bootrom.", false, true); err != nil {
			return errors.Wrap(err, "unable to register lpc flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveGenACPITables, "acpi", "A", true, "Generate ACPI tables. Required for FreeBSD/amd64 guests.", false, true); err != nil {
			return errors.Wrap(err, "unable to register GenACPITables flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveIncGuestCoreMem, "include-guest-mem", "C", false, "Include guest memory in core file.", false, true); err != nil {
			return errors.Wrap(err, "unable to register IncGuestCoreMem flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveExitOnUnemuIOPort, "exit-on-unemu-ioport", "e", false, "Force bhyve to exit when a guest issues an access to an I/O port that is not emulated.", false, true); err != nil {
			return errors.Wrap(err, "unable to register ExitOnUnemuIOPort flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveYieldCPUOnHLT, "yield-cpu-on-hlt", "H", true, "Yield the virtual CPU thread when a HLT instruction is detected.  If this option is not specified, virtual CPUs will use 100% of a host CPU.", false, true); err != nil {
			return errors.Wrap(err, "unable to register YieldCPUOnHLT flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveIgnoreUnimplementedMSRAccess, "ignore-unimp-msr-access", "w", false, "Ignore accesses to unimplemented Model Specific Registers.", false, true); err != nil {
			return errors.Wrap(err, "unable to register YieldCPUOnHLT flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveForceMSIInterrupts, "force-msii", "W", false, "Force virtio PCI device emulations to use MSI interrupts instead of MSI-X interrupts.", false, true); err != nil {
			return errors.Wrap(err, "unable to register ForceMSIInterrupts flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveExitOnPause, "exit-on-pause", "P", true, "Force the guest virtual CPU to exit when a PAUSE instruction is detected.", false, true); err != nil {
			return errors.Wrap(err, "unable to register ExitOnPause flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveWireGuestMemory, "wire-guest-memory", "S", true, "Wire guest memory.", false, true); err != nil {
			return errors.Wrap(err, "unable to register ExitOnPause flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveAPICx2Mode, "apicx2", "x", false, "The guest's local APIC is configured in x2APIC mode.", false, true); err != nil {
			return errors.Wrap(err, "unable to register APICx2Mode flag on VPC VM create")
		}

		if err := flag.AddBoolFlag(self, _KeyBhyveDisableMPTableGeneration, "disable-mpt-table-generation", "Y", false, "Disable MPtable generation.", false, true); err != nil {
			return errors.Wrap(err, "unable to register DisableMPTableGeneration flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyJSON, "json", "j", "", "JSON configuration file", false, false); err != nil {
			return errors.Wrap(err, "unable to register JSON flag on VPC Switch create")
		}

		if err := flag.AddStringFlag(self, _KeyName, "name", "n", "", "Virtual Machine Name", false, false); err != nil {
			return errors.Wrap(err, "unable to register Name flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyUUID, "uuid", "u", "", "VM UUID. A UUID is randomly generated by default.", true, false); err != nil {
			return errors.Wrap(err, "unable to register UUID flag on VPC VM create")
		}
		if err := flag.AddIntFlag(self, _KeyVCPUs, "vcpus", "", 1, "Number of vCPUs", false, false); err != nil {
			return errors.Wrap(err, "unable to register VCPUs flag on VPC VM create")
		}
		if err := flag.AddStringFlag(self, _KeyBootPartition, "bootpartition", "", "hd0,msdos1", "Partition to boot", false, false); err != nil {
			return errors.Wrap(err, "unable to register bootpartition flag on VPC VM create")
		}
		if err := flag.AddStringFlag(self, _KeyRAM, "ram", "", "256M", "RAM in MB e.g. 256M, 1G", false, false); err != nil {
			return errors.Wrap(err, "unable to register RAM flag on VPC VM create")
		}
		if err := flag.AddStringFlag(self, _KeyDiskDriver, "diskdriver", "", "virtio-blk", "Disk Driver Emulation", false, true); err != nil {
			return errors.Wrap(err, "unable to register DiskDriver flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyDiskDevice, "diskdevice", "", "", "Path to disk image/block device", false, true); err != nil {
			return errors.Wrap(err, "unable to register DiskDevice flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyDiskSize, "disksize", "", "256M", "Disk Size e.g. 256M, 10G", false, false); err != nil {
			return errors.Wrap(err, "unable to register DiskSize flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyNicDriver, "nicdriver", "", "kvirtio-net", "NIC Driver Emulation", false, true); err != nil {
			return errors.Wrap(err, "unable to register NicDriver flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyNicDevice, "nicdevice", "", "", "NIC Device Name e.g. vmnic0", false, false); err != nil {
			return errors.Wrap(err, "unable to register NicDevice flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeyNicID, "nicid", "", "", "NIC Device ID e.g. a vmnic uuid", false, false); err != nil {
			return errors.Wrap(err, "unable to register NicID flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeySerialConsole1, "serialconsole1", "", "", "Serial Console 1 Device", false, false); err != nil {
			return errors.Wrap(err, "unable to register SerialConsole1 flag on VPC VM create")
		}

		if err := flag.AddStringFlag(self, _KeySerialConsole2, "serialconsole2", "", "", "Serial Console 2 Device", false, false); err != nil {
			return errors.Wrap(err, "unable to register SerialConsole1 flag on VPC VM create")
		}
		return nil
	},
}
