//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

var Default = Build

var binary = "enbas"

// Test run the go tests.
// To enable verbose mode set ENBAS_TEST_VERBOSE=1.
// To enable coverage mode set ENBAS_TEST_COVER=1.
func Test() error {
	if err := changeToProjectRoot(); err != nil {
		return fmt.Errorf("unable to change to the project's root directory; %w", err)
	}

	goTest := sh.RunCmd("go", "test")

	args := []string{"./..."}

	if os.Getenv("ENBAS_TEST_VERBOSE") == "1" {
		args = append(args, "-v")
	}

	if os.Getenv("ENBAS_TEST_COVER") == "1" {
		args = append(args, "-cover")
	}

	return goTest(args...)
}

// Lint runs golangci-lint against the code.
func Lint() error {
	if err := changeToProjectRoot(); err != nil {
		return fmt.Errorf("unable to change to the project's root directory; %w", err)
	}

	return sh.RunV("golangci-lint", "run", "--color", "always")
}

// Build build the executable.
func Build() error {
	if err := changeToProjectRoot(); err != nil {
		return fmt.Errorf("unable to change to the project's root directory; %w", err)
	}

	main := "main.go"
	return sh.Run("go", "build", "-o", binary, main)
}

// Clean clean the workspace.
func Clean() error {
	if err := changeToProjectRoot(); err != nil {
		return fmt.Errorf("unable to change to the project's root directory; %w", err)
	}

	if err := sh.Rm(binary); err != nil {
		return err
	}

	if err := sh.Run("go", "clean", "./..."); err != nil {
		return err
	}

	return nil
}

func changeToProjectRoot() error {
	if err := os.Chdir("../.."); err != nil {
		return fmt.Errorf("unable to change directory; %w", err)
	}

	return nil
}
