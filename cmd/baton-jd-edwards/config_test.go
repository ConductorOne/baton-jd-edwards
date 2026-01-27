package main

import (
	"testing"

	cfg "github.com/conductorone/baton-jd-edwards/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/test"
	"github.com/conductorone/baton-sdk/pkg/ustrings"
)

func TestConfigs(t *testing.T) {
	test.ExerciseTestCasesFromExpressions(
		t,
		cfg.Config,
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
