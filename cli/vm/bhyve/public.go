package bhyve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	zfs "github.com/mistifyio/go-zfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/config"
	"github.com/y0ssar1an/q"
)

type BhyveConfig struct {
	Name                         string
	UUID                         string
	VCPUs                        int
	BootPartition                string
	RAM                          string
	DiskDriver                   string
	DiskDevice                   string
	DiskSize                     string
	NicDriver                    string
	NicDevice                    string
	NicID                        string
	SerialConsole1               string
	SerialConsole2               string
	LPC                          string
	HostBridge                   string
	GenACPITables                bool
	IncGuestCoreMem              bool
	ExitOnUnemuIOPort            bool
	ExitOnPause                  bool
	YieldCPUOnHLT                bool
	IgnoreUnimplementedMSRAccess bool
	ForceMSIInterrupts           bool
	Apicx2Mode                   bool
	DisableMPTableGeneration     bool
}

type BhyveArg interface {
	String() string
}

type BhyveFlag struct {
	flag string // -A, -H, -P, -S
}

func (flag BhyveFlag) String() string {
	return flag.flag
}

type BhyveKV struct {
	Key   string // i.e. -m for memory, or -c for cpu
	Value string // i.e. 32G for memory, 16 for cpu
}

type BhyveKVs []BhyveKV

func (kv BhyveKV) String() string {
	if kv.Value == "" {
		return kv.Key
	}
	return strings.Join([]string{kv.Key, kv.Value}, "=")
}

func (kvs BhyveKVs) String() string {
	var s []string
	for _, kv := range kvs {
		s = append(s, kv.String())
	}
	return strings.Join(s, " ")
}

// BhyveDevice defines a virtual PCI slot and function
// -s slot,emulation[,conf]
type BhyveDevice struct {
	Flag      string // i.e. -s or -l
	Arg       string
	Slot      string
	Emulation string // i.e. hostbridge, passthru, virtio-net, virtio-blk, lpc
	Conf      BhyveKV
}

// TODO(SAM): Clean up the comma algorithm
func (bd BhyveDevice) String() string {
	var str string
	args := []string{bd.Slot, bd.Emulation, bd.Conf.String()}
	active := []string{}
	for _, arg := range args {
		if arg != "" {
			active = append(active, arg)
		}
	}
	q.Q(active)
	if bd.Arg != "" {
		str += bd.Arg
	}
	for i, arg := range active {
		str += fmt.Sprintf("%s", arg)
		if i != len(active) && i+1 < len(active) {
			str += ","
		}
	}
	return strings.Join([]string{bd.Flag, str}, " ")
}

// TODO(SAM): Clean up the comma algorithm
func (bd BhyveDevice) Array() []string {
	var str string
	args := []string{bd.Slot, bd.Emulation, bd.Conf.String()}
	active := []string{}
	for _, arg := range args {
		if arg != "" {
			active = append(active, arg)
		}
	}
	q.Q(active)
	if bd.Arg != "" {
		str += bd.Arg
	}
	for i, arg := range active {
		str += fmt.Sprintf("%s", arg)
		if i != len(active) && i+1 < len(active) {
			str += ","
		}
	}
	if bd.Flag != "" {
		return []string{bd.Flag, str}
	}
	return []string{str}
}

type BhyveDeviceList []BhyveDevice

func (bdl BhyveDeviceList) String() string {
	var s []string
	for _, bd := range bdl {
		if arg := bd.String(); arg != "" {
			s = append(s, arg)
		}
	}
	return strings.Join(s, " ")
}

func (bdl BhyveDeviceList) Array() []string {
	var s []string
	for _, bd := range bdl {
		for _, arg := range bd.Array() {
			if arg != "" {

				s = append(s, arg)
			}
		}
	}
	return s
}

type BhyveCommand struct {
	BinaryPath string      // Path to bhyve e.g. /usr/sbin/bhyve
	Flags      []BhyveFlag // Boolean Flags such as -A -H -P -S
	Args       []BhyveArg  // All arguments with a value associated e.g. -m 32G
	Devices    BhyveDeviceList
	Name       string // UUID of VM
}

func PrintConfig(config BhyveConfig) {
	log.Info().
		Str("name", config.Name).
		Str("uuid", config.UUID).
		Int("vcpus", config.VCPUs).
		Str("bootpartition", config.BootPartition).
		Str("ram", config.RAM).
		Str("diskdriver", config.DiskDriver).
		Str("diskdevice", config.DiskDevice).
		Str("disksize", config.DiskSize).
		Str("nicdriver", config.NicDriver).
		Str("nicdevice", config.NicDevice).
		Str("serialconsole1", config.SerialConsole1).
		Str("serialconsole2", config.SerialConsole2).
		Str("hostbridge", config.HostBridge).
		Str("lpc", config.LPC).
		Bool("genacpitables", config.GenACPITables).
		Bool("incguestcoremem", config.IncGuestCoreMem).
		Bool("exitonunemuioport", config.ExitOnUnemuIOPort).
		Bool("yieldcpuonhlt", config.YieldCPUOnHLT).
		Bool("ignoreunimplementedmsraccess", config.IgnoreUnimplementedMSRAccess).
		Bool("forcemsiinterrupts", config.ForceMSIInterrupts).
		Bool("apicx2mode", config.Apicx2Mode).
		Bool("disablemptablegeneration", config.DisableMPTableGeneration).
		Msg("Bhyve Configuration")
}

func buildGrubBhyveArg(cfg BhyveConfig, arg string) (*BhyveDevice, error) {
	switch arg {
	case config.KeyBhyveWireGuestMemory:
		return &BhyveDevice{Flag: "-S"}, nil
	case config.KeyDiskDevice:
		deviceMapPath, err := GetDeviceMap(cfg.UUID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get device map path")
		}
		return &BhyveDevice{Flag: "-m", Arg: deviceMapPath}, nil
	case config.KeyBootPartition:
		return &BhyveDevice{Flag: "-r", Arg: cfg.BootPartition}, nil
	case config.KeySerialConsole1:
		return &BhyveDevice{Flag: "-c", Arg: cfg.SerialConsole1}, nil
	case config.KeyRAM:
		return &BhyveDevice{Flag: "-M", Arg: cfg.RAM}, nil
	case config.KeyName:
		return &BhyveDevice{Arg: cfg.Name}, nil
	case config.KeyUUID:
		// Using only the first part of the UUID, VM names cannot contain an hypen
		// TODO(Sam) Are we sure we won't run into collisions here?
		uuid := strings.Split(cfg.UUID, "-")[0]
		return &BhyveDevice{Arg: uuid}, nil
	}
	return nil, errors.Errorf("no such flag: %s", arg)
}

func buildBhyveArg(cfg BhyveConfig, arg string) (*BhyveDevice, error) {
	switch arg {
	case config.KeyBhyveGenACPITables:
		return &BhyveDevice{Flag: "-A"}, nil
	case config.KeyBhyveYieldCPUOnHLT:
		return &BhyveDevice{Flag: "-H"}, nil
	case config.KeyBhyveExitOnUnemuIOPort:
		return &BhyveDevice{Flag: "-e"}, nil
	case config.KeyBhyveExitOnPause:
		return &BhyveDevice{Flag: "-P"}, nil
	case config.KeyBhyveWireGuestMemory:
		return &BhyveDevice{Flag: "-S"}, nil
	case config.KeyDiskDevice:
		return &BhyveDevice{Flag: "-s", Slot: "4", Emulation: cfg.DiskDriver, Conf: BhyveKV{Key: cfg.DiskDevice}}, nil
	case config.KeyBootPartition:
		return &BhyveDevice{Flag: "-r", Arg: cfg.BootPartition}, nil
	case config.KeySerialConsole1:
		return &BhyveDevice{Flag: "-l", Emulation: "com1", Conf: BhyveKV{Key: cfg.SerialConsole1}}, nil
	case config.KeyVCPUs:
		return &BhyveDevice{Flag: "-c", Arg: strconv.Itoa(cfg.VCPUs)}, nil
	case config.KeyRAM:
		return &BhyveDevice{Flag: "-m", Arg: cfg.RAM}, nil
	case config.KeyNicDevice:
		return &BhyveDevice{Flag: "-s", Slot: "5", Emulation: cfg.NicDriver, Conf: BhyveKV{Key: "id", Value: cfg.NicID}}, nil
	case config.KeyHostBridge:
		return &BhyveDevice{Flag: "-s", Slot: "0", Emulation: cfg.HostBridge}, nil
	case config.KeyLPC:
		return &BhyveDevice{Flag: "-s", Slot: "31", Emulation: cfg.LPC}, nil
	case config.KeyName:
		return &BhyveDevice{Arg: cfg.Name}, nil
	case config.KeyUUID:
		// Using only the first part of the UUID, VM names cannot contain an hypen
		// TODO(Sam) Are we sure we won't run into collisions here?
		uuid := strings.Split(cfg.UUID, "-")[0]
		return &BhyveDevice{Arg: uuid}, nil
	}
	return nil, errors.Errorf("no such flag: %s", arg)
}

func BuildGrubBhyveArgs(cfg BhyveConfig) (*BhyveCommand, error) {
	grubBhyveRelevantArgs := []string{config.KeyBhyveWireGuestMemory, config.KeyDiskDevice, config.KeyBootPartition, config.KeySerialConsole1, config.KeyRAM, config.KeyUUID}
	grubBhyvePath := "/usr/local/sbin/grub-bhyve"

	var args []BhyveDevice
	for _, a := range grubBhyveRelevantArgs {
		arg, err := buildGrubBhyveArg(cfg, a)
		if err != nil {
			return nil, errors.Wrap(err, "Device is nil")
		}
		args = append(args, *arg)
	}

	return &BhyveCommand{BinaryPath: grubBhyvePath, Devices: args}, nil
}

func RunGrubBhyve(cfg BhyveConfig, cmd *BhyveCommand) error {
	var outb, errb bytes.Buffer
	//grubBhyveString := BuildGrubBhyveArgs(cfg)
	c := exec.Command(cmd.BinaryPath, cmd.Devices.Array()...)
	c.Stdout = &outb
	c.Stderr = &errb
	q.Q(c)
	q.Q(cmd.Devices.Array())
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "grub-bhyve failed")
	}

	fmt.Printf("stdout: %s\n", outb.String())
	fmt.Printf("stderr: %s\n", errb.String())
	log.Info().
		Str("command", cmd.BinaryPath).
		Str("args", cmd.Devices.String()).
		Msg("Running grub-bhyve")
	return nil
}

func BuildBhyveArgs(cfg BhyveConfig) (*BhyveCommand, error) {
	// bhyveRelevantArgs := []string{config.KeyBhyveGenACPITables, config.KeyBhyveYieldCPUOnHLT, config.KeyBhyveExitOnPause, config.KeyBhyveWireGuestMemory, config.KeyUUID, config.KeyRAM, config.KeyDiskDevice, config.KeyNicDevice, config.KeySerialConsole1, config.KeyName}
	bhyveRelevantArgs := []string{config.KeyBhyveGenACPITables, config.KeyBhyveYieldCPUOnHLT, config.KeyBhyveExitOnPause, config.KeyBhyveWireGuestMemory, config.KeyVCPUs, config.KeyRAM, config.KeyDiskDevice, config.KeySerialConsole1, config.KeyName, config.KeyHostBridge, config.KeyLPC, config.KeyUUID}
	bhyvePath := "/usr/sbin/bhyve"

	var args []BhyveDevice
	for _, a := range bhyveRelevantArgs {
		arg, err := buildBhyveArg(cfg, a)
		q.Q(arg)
		if err != nil {
			return nil, errors.Wrap(err, "Device is nil")
		}
		if arg != nil {
			args = append(args, *arg)
		}
	}

	return &BhyveCommand{BinaryPath: bhyvePath, Devices: args}, nil
}

func RunBhyve(cfg BhyveConfig, cmd *BhyveCommand) error {
	var outb, errb bytes.Buffer
	//grubBhyveString := BuildGrubBhyveArgs(cfg)
	c := exec.Command(cmd.BinaryPath, cmd.Devices.Array()...)
	c.Stdout = &outb
	c.Stderr = &errb
	q.Q(c)
	q.Q(cmd.Devices.Array())
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "bhyve failed")
	}

	fmt.Printf("stdout: %s\n", outb.String())
	fmt.Printf("stderr: %s\n", errb.String())
	log.Info().
		Str("command", cmd.BinaryPath).
		Str("args", cmd.Devices.String()).
		Msg("Running bhyve")
	return nil
}

func WriteConfig(config BhyveConfig) error {
	configJson, err := json.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "unable to marshal bhyve config into json")
	}

	guestPath, err := GetGuestPath(config.UUID)
	if err != nil {
		return errors.Wrap(err, "unable to get guest path")
	}

	configPath := fmt.Sprintf("%s/config.json", guestPath)
	err = ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		return errors.Wrap(err, "unable to write config")
	}

	log.Info().
		Str("configPath", configPath).
		Msg("Wrote VM Config")
	return nil
}

func ReadConfig(uuid string) (*BhyveConfig, error) {
	var configJson *BhyveConfig

	guestPath, err := GetGuestPath(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get guest path")
	}

	configPath := fmt.Sprintf("%s/config.json", guestPath)
	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read config")
	}

	json.Unmarshal(raw, &configJson)

	log.Info().
		Str("configPath", configPath).
		Msg("Read VM Config")
	return configJson, nil
}

func GetZpoolName() (string, error) {
	zpools, err := zfs.ListZpools()
	if err != nil {
		return "", errors.Wrap(err, "unable to find zpool")
	}

	// Always returns the first zpool found
	zpoolName := zpools[0].Name
	log.Info().Str("zpool", zpoolName).Msg("Using zpool")

	return zpoolName, nil
}

func GetGuestDataset(uuid string) (string, error) {
	zpoolName, err := GetZpoolName()
	if err != nil {
		return "", errors.Wrap(err, "unable to get zpool name")
	}

	path := fmt.Sprintf("%s/guests/%s", zpoolName, uuid)

	return path, nil
}

// To create the guest path string we merely need to prefix the dataset with '/'
func GetGuestPath(uuid string) (string, error) {
	datasetPath, err := GetGuestDataset(uuid)
	if err != nil {
		return "", errors.Wrap(err, "unable to get guest dataset")
	}

	path := fmt.Sprintf("/%s", datasetPath)

	return path, nil
}

func GetDiskPath(uuid string) (string, error) {
	guestPath, err := GetGuestPath(uuid)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/disk0", guestPath)
	return path, nil
}

func GetZvolPath(uuid string) (string, error) {
	guestPath, err := GetDiskPath(uuid)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/dev/zvol%s", guestPath)
	return path, nil
}

func GetDeviceMap(uuid string) (string, error) {
	guestPath, err := GetGuestPath(uuid)
	if err != nil {
		return "", errors.Wrap(err, "unable to get guest path")
	}

	return fmt.Sprintf("%s/device.map", guestPath), nil
}
