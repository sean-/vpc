---
date: 2018-02-28T23:59:59Z
title: "vpc switch port add"
slug: vpc_switch_port_add
url: /command/vpc_switch_port_add
---
## vpc switch port add

add a port to a VPC Switch

### Synopsis


add a port to a VPC Switch

```
vpc switch port add [flags]
```

### Options

```
      --ethlink-id string   Specify the ID of the VPC EthLink
  -h, --help                help for add
  -n, --l2-name string      Name of the underlying L2 interface to be wrapped as a VPC EthLink and used as the uplink in the VPC Switch
      --port-id string      Specify the VPC Port ID
      --switch-id string    Specify the VPC Switch ID
  -u, --uplink              make the port ID an uplink for the switch
```

### Options inherited from parent commands

```
  -F, --log-format string   Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string    Change the log level being sent to stdout (default "INFO")
      --use-color           Use ASCII colors
  -P, --use-pager           Use a pager to read the output (defaults to $PAGER, less(1), or more(1))
  -Z, --utc                 Display times in UTC
```

### SEE ALSO
* [vpc switch port](/command/vpc_switch_port)	 - VPC switch management

