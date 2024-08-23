package main

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/test"
	"github.com/conductorone/baton-sdk/pkg/ustrings"
)

func TestConfigs(t *testing.T) {
	test.ExerciseTestCasesFromExpressions(
		t,
		field.NewConfiguration(configurationFields),
		nil,
		ustrings.ParseFlags,
		[]test.TestCaseFromExpression{
			{
				"",
				false,
				"missing required fields",
			},
			{
				"--ais-url 1 --username 1 --password 1",
				true,
				"is valid",
			},
			{
				"--ais-url 1 --username 1 --password 1 --env 1",
				true,
				"is valid with optional field",
			},
		},
	)
}
