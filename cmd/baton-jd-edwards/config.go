package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	aisUrlField = field.StringField(
		"ais-url",
		field.WithRequired(true),
		field.WithDescription("Your JD Edwards AIS Server REST API url. Provided url should contain port. (e.g: https://your_ais_server:port)."),
	)
	usernameField = field.StringField(
		"username",
		field.WithRequired(true),
		field.WithDescription("JD Edwards EnterpriseOne username."),
	)
	passwordField = field.StringField(
		"password",
		field.WithRequired(true),
		field.WithDescription("JD Edwards EnterpriseOne password."),
	)
	envField = field.StringField(
		"env",
		field.WithDescription("Environment to use for login. If not specified, the default environment configured for the AIS Server will be used."),
	)
	configurationFields = []field.SchemaField{
		aisUrlField,
		usernameField,
		passwordField,
		envField,
	}
)
