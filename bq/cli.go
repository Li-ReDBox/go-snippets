package bq

import (
	"log"
	"os"
)

// Required checks command-line flag which also can be sourced from environment variable
// If both are empty, it is an fatal error
func Required(input *string, envName, name string) {
	if *input != "" {
		return
	}
	*input = os.Getenv(envName)
	if *input == "" {
		log.Fatalf("%s is neither set on command line or by environment variable %s.\n", name, envName)
	}
}
