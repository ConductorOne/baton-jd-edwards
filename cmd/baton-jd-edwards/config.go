package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	AisUrl   string `mapstructure:"ais-url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.AisUrl == "" {
		return fmt.Errorf("AIS URL is missing, please provide it via --ais-url flag or $BATON_AIS_URL environment variable")
	}

	if cfg.Username == "" {
		return fmt.Errorf("username is missing, please provide it via --username flag or $BATON_USERNAME environment variable")
	}

	if cfg.Password == "" {
		return fmt.Errorf("password is missing, please provide it via --password flag or $BATON_PASSWORD environment variable")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("ais-url", "", "Your JD Edwards Application Interface Service (AIS) Server REST API url. Provided url should contain port. (e.g: https://your_ais_server:port)($BATON_AIS_URL)")
	cmd.PersistentFlags().String("username", "", "JD Edwards EnterpriseOne username. ($BATON_USERNAME)")
	cmd.PersistentFlags().String("password", "", "JD Edwards EnterpriseOne password. ($BATON_PASSWORD)")
}
