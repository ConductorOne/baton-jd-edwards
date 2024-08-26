package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-jd-edwards/pkg/connector"
	configSchema "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	version       = "dev"
	connectorName = "baton-jd-edwards"
)

func main() {
	ctx := context.Background()
	_, cmd, err := configSchema.DefineConfiguration(
		ctx,
		connectorName,
		getConnector,
		field.NewConfiguration(configurationFields),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version
	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	cb, err := connector.New(ctx,
		cfg.GetString(aisUrlField.FieldName),
		cfg.GetString(usernameField.FieldName),
		cfg.GetString(passwordField.FieldName),
		cfg.GetString(envField.FieldName),
	)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return c, nil
}
