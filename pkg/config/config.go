package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AisUrlField = field.StringField(
		"ais-url",
		field.WithRequired(true),
		field.WithDescription("Your JD Edwards AIS Server REST API url. Provided url should contain port. (e.g: https://your_ais_server:port)."),
	)
	UsernameField = field.StringField(
		"username",
		field.WithRequired(true),
		field.WithDescription("JD Edwards EnterpriseOne username."),
	)
	PasswordField = field.StringField(
		"password",
		field.WithRequired(true),
		field.WithDescription("JD Edwards EnterpriseOne password."),
	)
	EnvField = field.StringField(
		"env",
		field.WithDescription("Environment to use for login. If not specified, the default environment configured for the AIS Server will be used."),
	)

	ConfigurationFields = []field.SchemaField{
		AisUrlField,
		UsernameField,
		PasswordField,
		EnvField,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields)
