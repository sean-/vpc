package agent

import (
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/cmd/vpc/config"
	"github.com/sean-/vpc/db"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "agent"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "Run " + buildtime.PROGNAME,
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "run").Msg("")

			// 1. Parse config and construct agent
			var config config.Config
			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to decode config into struct")
			}

			dbPool, err := db.New(config.DBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create database pool")
			}

			// 2. Run agent
			a, err := agent.New()
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}

			a.DbPool = dbPool

			if err := a.Start(); err != nil {
				return errors.Wrapf(err, "unable to start %s agent", buildtime.PROGNAME)
			}
			defer a.Stop()

			// 3. Connect to the database to verify database credentials
			if err := a.DbPool.Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}

			// 4. Loop until program exit
			if err := a.Run(); err != nil {
				return errors.Wrapf(err, "unable to run %s agent", buildtime.PROGNAME)
			}

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		viper.SetDefault("db.scheme", "crdb")
		viper.SetDefault("db.user", "root")
		viper.SetDefault("db.host", "localhost")
		viper.SetDefault("db.port", 26257)
		viper.SetDefault("db.database", "triton")
		viper.SetDefault("db.conn_timeout", 10*time.Second)
		viper.SetDefault("db.insecure_skip_verify", false)

		caPath, err := homedir.Expand("~/.cockroach-certs/ca.crt")
		if err != nil {
			return errors.Wrap(err, "error expanding home directory")
		}
		viper.SetDefault("db.ca_path", caPath)

		certPath, err := homedir.Expand("~/.cockroach-certs/client.root.crt")
		if err != nil {
			return errors.Wrap(err, "error expanding home directory")
		}
		viper.SetDefault("db.cert_path", certPath)

		keyPath, err := homedir.Expand("~/.cockroach-certs/client.root.key")
		if err != nil {
			return errors.Wrap(err, "error expanding home directory")
		}
		viper.SetDefault("db.key_path", keyPath)

		return nil
	},
}
