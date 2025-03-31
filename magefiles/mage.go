//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	projectName          = "enbas"
	defaultAppName       = projectName
	defaultInstallPrefix = "/usr/local"
	toolsModFile         = "./tools/tools.mod"

	envInstallPrefix    = "ENBAS_INSTALL_PREFIX"
	envTestVerbose      = "ENBAS_TEST_VERBOSE"
	envTestCover        = "ENBAS_TEST_COVER"
	envBuildRebuildAll  = "ENBAS_BUILD_REBUILD_ALL"
	envBuildVerbose     = "ENBAS_BUILD_VERBOSE"
	envFailOnFormatting = "ENBAS_FAIL_ON_FORMATTING"
	envAppName          = "ENBAS_APP_NAME"
	envAppVersion       = "ENBAS_APP_VERSION"
)

var Default = Build

// Test run the go tests.
// To enable verbose mode set ENBAS_TEST_VERBOSE=1.
// To enable coverage mode set ENBAS_TEST_COVER=1.
func Test() error {
	goTest := sh.RunCmd("go", "test")

	args := []string{"./..."}

	if os.Getenv(envTestVerbose) == "1" {
		args = append(args, "-v")
	}

	if os.Getenv(envTestCover) == "1" {
		args = append(args, "-cover")
	}

	return goTest(args...)
}

// Lint runs golangci-lint against the code.
func Lint() error {
	return sh.RunV("golangci-lint", "run", "--color", "always")
}

// Gosec runs gosec against the code.
func Gosec() error {
	return sh.RunV(
		"go",
		"tool",
		"--modfile",
		toolsModFile,
		"gosec",
		"./...",
	)
}

// Staticcheck runs staticcheck against the code.
func Staticcheck() error {
	return sh.RunV(
		"go",
		"tool",
		"--modfile",
		toolsModFile,
		"staticcheck",
		"./...",
	)
}

// Gofmt checks the code for formatting.
// To fail on formatting set BEACON_FAIL_ON_FORMATTING=1
func Gofmt() error {
	output, err := sh.Output("go", "fmt", "./...")
	if err != nil {
		return err
	}

	formattedFiles := ""

	for file := range strings.SplitSeq(output, "\n") {
		formattedFiles += "\n- " + file
	}

	if os.Getenv(envFailOnFormatting) != "1" {
		fmt.Println(formattedFiles)

		return nil
	}

	if len(output) != 0 {
		return fmt.Errorf("the following files needed to be formatted: %s", formattedFiles)
	}

	return nil
}

// Govet runs go vet against the code.
func Govet() error {
	return sh.RunV("go", "vet", "./...")
}

// Build build the executable.
// To rebuild packages that are already up-to-date set ENBAS_BUILD_REBUILD_ALL=1
// To enable verbose mode set ENBAS_BUILD_VERBOSE=1
func Build() error {
	fmt.Println("Building the binary...")

	main := "./cmd/" + projectName
	flags := ldflags()
	build := sh.RunCmd("go", "build")
	args := []string{"-ldflags=" + flags, "-o", binary()}

	if os.Getenv(envBuildRebuildAll) == "1" {
		args = append(args, "-a")
	}

	if os.Getenv(envBuildVerbose) == "1" {
		args = append(args, "-v")
	}

	args = append(args, main)

	if err := build(args...); err != nil {
		return fmt.Errorf("error building the binary: %w", err)
	}

	return nil
}

// Install install the executable.
func Install() error {
	mg.Deps(Build)

	installPrefix := os.Getenv(envInstallPrefix)
	app := appName()
	binary := binary()

	if installPrefix == "" {
		installPrefix = defaultInstallPrefix
	}

	dest := filepath.Join(installPrefix, "bin", app)

	fmt.Println("Installing the binary to", dest)

	if err := sh.Copy(dest, binary); err != nil {
		return fmt.Errorf("unable to install %s; %w", dest, err)
	}

	fmt.Printf("Successfully installed %s to %s\n", projectName, dest)

	return nil
}

// Clean clean the workspace.
func Clean() error {
	fmt.Println("Cleaning the workspace...")

	if err := sh.Rm(binary()); err != nil {
		return err
	}

	if err := sh.Run("go", "clean", "./..."); err != nil {
		return err
	}

	fmt.Println("Workspace cleaned.")

	return nil
}

// ldflags returns the build flags.
func ldflags() string {
	var (
		infoPackage              = "codeflow.dananglin.me.uk/apollo/enbas/internal/info"
		binaryVersionVar         = infoPackage + "." + "BinaryVersion"
		gitCommitVar             = infoPackage + "." + "GitCommit"
		buildTimeVar             = infoPackage + "." + "BuildTime"
		applicationNameVar       = infoPackage + "." + "ApplicationName"
		applicationTitledNameVar = infoPackage + "." + "ApplicationTitledName"

		ldflagsfmt = "-s -w -X %s=%s -X %s=%s -X %s=%s -X %s=%s -X %s=%s"
		buildTime  = time.Now().UTC().Format(time.RFC3339)
	)

	return fmt.Sprintf(
		ldflagsfmt,
		binaryVersionVar, binaryVersion(),
		gitCommitVar, gitCommit(),
		buildTimeVar, buildTime,
		applicationNameVar, appName(),
		applicationTitledNameVar, appTitledName(),
	)
}

// binaryVersion returns the version of the binary.
// If ENBAS_APP_VERSION is set, the value of that is returned, otherwise
// the latest git tag (using git describe) is returned.
func binaryVersion() string {
	ver := os.Getenv(envAppVersion)
	if ver != "" {
		return ver
	}

	ver, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		fmt.Printf("WARNING: error getting the binary version: %v.\n", err)
		return "N/A"
	}

	return ver
}

// gitCommit returns the current git commit
func gitCommit() string {
	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		commit = "N/A"
	}

	return commit
}

// appName returns the application's name.
// The value of ENBAS_APP_NAME is return if the environment variable is set
// otherwise the default name is returned.
func appName() string {
	appName := os.Getenv(envAppName)

	if appName == "" {
		return defaultAppName
	}

	return appName
}

func binary() string {
	return filepath.Join("./__build", appName())
}

func appTitledName() string {
	runes := []rune(appName())
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
