//go:build mage

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

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

	flags := ldflags()
	return sh.Run("go", "build", "-ldflags="+flags, "-a", "-o", binary, "./cmd/enbas")
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

// ldflags returns the build flags.
func ldflags() string {
	ldflagsfmt := "-s -w -X main.binaryVersion=%s -X main.gitCommit=%s -X main.goVersion=%s -X main.buildTime=%s"
	buildTime := time.Now().UTC().Format(time.RFC3339)

	return fmt.Sprintf(ldflagsfmt, version(), gitCommit(), runtime.Version(), buildTime)
}

// version returns the latest git tag using git describe.
func version() string {
	version, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		version = "N/A"
	}

	return version
}

// gitCommit returns the current git commit
func gitCommit() string {
	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		commit = "N/A"
	}

	return commit
}
