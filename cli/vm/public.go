package vm

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/cli/vm/create"
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
		//cmds := []*command.Command{
		//	create.Cmd,
		//	// list.Cmd,
		//}

		//for _, cmd := range cmds {
		//	parent.Cobra.AddCommand(cmd.Cobra)
		//	cmd.Setup(cmd)
		//}
		subCommands := command.Commands{
			create.Cmd,
			//list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
