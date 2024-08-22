package main

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/test"
	"github.com/conductorone/baton-sdk/pkg/ustrings"
	"github.com/spf13/viper"
)

type TestCaseFromExpression = struct {
	Expression string
	IsValid    bool
	Message    string
}

func ExerciseTestCase(
	t *testing.T,
	configurationSchema field.Configuration,
	extraValidationFunction func(*viper.Viper) error,
	configs map[string]string,
	isValid bool,
) {
	test.AssertValidation(
		t,
		func() error {
			v := test.MakeViper(configs)
			err := field.Validate(configurationSchema, v)
			if err != nil {
				return err
			}
			if extraValidationFunction != nil {
				return extraValidationFunction(v)
			}
			return nil
		},
		isValid,
	)
}

// ExerciseTestCasesFromExpressions - Like ExerciseTestCases, but instead of
// passing a `map[string]string` to each test case, pass a function that parses
// configs from strings and pass each test case an expression as a string.
func ExerciseTestCasesFromExpressions(
	t *testing.T,
	configurationSchema field.Configuration,
	extraValidationFunction func(*viper.Viper) error,
	expressionParser func(string) (map[string]string, error),
	testCases []TestCaseFromExpression,
) {
	for _, testCase := range testCases {
		t.Run(testCase.Message, func(t *testing.T) {
			values, err := expressionParser(testCase.Expression)
			if err != nil {
				t.Fatal("could not parse flags:", err)
			}
			ExerciseTestCase(
				t,
				configurationSchema,
				extraValidationFunction,
				values,
				testCase.IsValid,
			)
		})
	}
}

func TestConfigs(t *testing.T) {
	ExerciseTestCasesFromExpressions(
		t,
		field.NewConfiguration(configurationFields),
		nil,
		ustrings.ParseFlags,
		[]TestCaseFromExpression{
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
