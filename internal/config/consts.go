package config

const (
	DefaultManDir = "./docs/man"
	ManSect       = 8

	DefaultMarkdownDir       = "./docs/md"
	DefaultMarkdownURLPrefix = "/command"

	KeyDocManDir            = "doc.mandir"
	KeyDocMarkdownDir       = "doc.markdown-dir"
	KeyDocMarkdownURLPrefix = "doc.markdown-url-prefix"

	KeyEthLinkDestroyID  = "ethlink.destroy.ethlink-id"
	KeyEthLinkListSortBy = "ethlink.list.sort-by"

	KeyListObjCounts = "list.obj-counts"
	KeyListObjSortBy = "list.sort-by"
	KeyListObjType   = "list.type"

	KeyLogFormat    = "log.format"
	KeyLogLevel     = "log.level"
	KeyLogStats     = "log.stats"
	KeyLogTermColor = "log.use-color"

	KeySWPortAddEthLinkID          = "switch.port.add.ethlink-id"
	KeySWPortAddID                 = "switch.port.add.id"
	KeySWPortAddL2Name             = "switch.port.add.l2-name"
	KeySWPortAddMAC                = "switch.port.add.mac"
	KeySWPortAddSwitchID           = "switch.port.add.switch-id"
	KeySWPortAddUplink             = "switch.port.add.uplink"
	KeySWPortConnectInterfaceID    = "switch.port.connect.interface-id"
	KeySWPortConnectPortID         = "switch.port.connect.port-id"
	KeySWPortDisconnectInterfaceID = "switch.port.disconnect.interface-id"
	KeySWPortDisconnectPortID      = "switch.port.disconnect.port-id"

	KeySWPortRemovePortID   = "switch.port.remove.port-id"
	KeySWPortRemoveSwitchID = "switch.port.remove.switch-id"

	KeyShellAutoCompBashDir = "shell.autocomplete.bash-dir"

	KeySWCreateSwitchID  = "switch.create.switch-id"
	KeySWCreateSwitchMAC = "switch.create.switch-mac"
	KeySWCreateVNI       = "switch.create.vni"
	KeySWDestroySwitchID = "switch.destroy.switch-id"

	KeyUsePager = "general.use-pager"
	KeyUseUTC   = "general.utc"

	KeyVMNICCreateID    = "vmnic.create.id"
	KeyVMNICCreateMAC   = "vmnic.create.mac"
	KeyVMNICDestroyID   = "vmnic.destroy.id"
	KeyVMNICGetNQueues  = "vmnic.get.num-queues"
	KeyVMNICGetVMNICID  = "vmnic.get.vmnic-id"
	KeyVMNICSetFreeze   = "vmnic.set.freeze"
	KeyVMNICSetNQueues  = "vmnic.set.num-queues"
	KeyVMNICSetUnfreeze = "vmnic.set.unfreeze"
	KeyVMNICSetVMNICID  = "vmnic.set.vmnic-id"

	KeyBhyveGenACPITables                = "genacpitables"
	KeyBhyveIncGuestCoreMem              = "incguestcoremem"
	KeyBhyveExitOnUnemuIOPort            = "exitonunemuioport"
	KeyBhyveYieldCPUOnHLT                = "yieldcpuonhlt"
	KeyBhyveIgnoreUnimplementedMSRAccess = "ignoreunimplementedmsraccess"
	KeyBhyveForceMSIInterrupts           = "forcemsiinterrupts"
	KeyBhyveAPICx2Mode                   = "apicx2mode"
	KeyBhyveDisableMPTableGeneration     = "disablemptablegeneration"
	KeyBhyveExitOnPause                  = "exitonpause"
	KeyBhyveWireGuestMemory              = "wireguestmemory"

	KeyJson           = "json"
	KeyName           = "vmname"
	KeyUUID           = "uuid"
	KeyVCPUs          = "vcpus"
	KeyBootPartition  = "bootpartition"
	KeyRAM            = "ram"
	KeyDiskDriver     = "diskdriver"
	KeyDiskDevice     = "diskdevice"
	KeyDiskSize       = "disksize"
	KeyNicDriver      = "nicdriver"
	KeyNicDevice      = "nicdevice"
	KeyNicID          = "nicid"
	KeySerialConsole1 = "serialconsole1"
	KeySerialConsole2 = "serialconsole2"
	KeyHostBridge     = "hostbridge"
	KeyLPC            = "lpc"
)
