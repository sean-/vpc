package shell

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/vpc/shell/autocompletion"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const cmdName = "shell"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:   cmdName,
		Short: "shell commands",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			autocompletion.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", cmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
