package bq

import (
	"os"
	"testing"
)

func Test_required(t *testing.T) {
	envName := "ENV_VAR"

	t.Run("input has the value", func(t *testing.T) {
		input := "set"

		os.Setenv("test", "replaced")

		Required(&input, envName, "input has been set")
		defer os.Unsetenv(envName)

		if input != "set" {
			t.Errorf("input set has been changed to %s", input)
		}
	})

	t.Run("environment variable is the source", func(t *testing.T) {
		empty := ""

		os.Setenv(envName, "replaced")
		defer os.Unsetenv(envName)

		Required(&empty, envName, "env is the source")

		if empty != "replaced" {
			t.Error("environment variable does not replase an empty input")
		}
	})
}
