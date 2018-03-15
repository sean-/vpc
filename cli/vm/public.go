package vm

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/cli/vm/create"
	"github.com/sean-/vpc/cli/vm/list"
	vmstart "github.com/sean-/vpc/cli/vm/start"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "vm"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:   _CmdName,
		Short: "VM management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			create.Cmd,
			list.Cmd,
			vmstart.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
