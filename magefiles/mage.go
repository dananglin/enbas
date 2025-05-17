//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	projectName    = "enbas"
	defaultAppName = projectName
	toolsModFile   = "./tools/tools.mod"
	buildDir       = "./__build"
	binDir         = buildDir + "/bin"
	manDir         = buildDir + "/share/man"

	envTestVerbose      = "ENBAS_TEST_VERBOSE"
	envTestCover        = "ENBAS_TEST_COVER"
	envBuildRebuildAll  = "ENBAS_BUILD_REBUILD_ALL"
	envBuildVerbose     = "ENBAS_BUILD_VERBOSE"
	envFailOnFormatting = "ENBAS_FAIL_ON_FORMATTING"
	envAppName          = "ENBAS_APP_NAME"
	envAppVersion       = "ENBAS_APP_VERSION"
	envAppCommitRef     = "ENBAS_APP_COMMIT_REF"
)

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

type Build mg.Namespace

// Binary build the binary.
// To rebuild packages that are already up-to-date set ENBAS_BUILD_REBUILD_ALL=1
// To enable verbose mode set ENBAS_BUILD_VERBOSE=1
func (Build) Binary() error {
	fmt.Println("Building the binary...")

	if err := ensureDirectory(binDir); err != nil {
		return fmt.Errorf("error ensuring the presence of %q: %w", binDir, err)
	}

	build := sh.RunCmd("go", "build")
	args := []string{"-ldflags=" + ldflags(), "-o", binary()}

	if os.Getenv(envBuildRebuildAll) == "1" {
		args = append(args, "-a")
	}

	if os.Getenv(envBuildVerbose) == "1" {
		args = append(args, "-v")
	}

	args = append(args, "./cmd/"+projectName)

	if err := build(args...); err != nil {
		return fmt.Errorf("error building the binary: %w", err)
	}

	fmt.Println("Binary successfully built to ", binary())

	return nil
}

// Documentation builds the man pages and the example configuration.
func (Build) Documentation() error {
	fmt.Println("Generating the documentation...")

	examplesDir := buildDir + "/share/doc/" + appName() + "/examples"

	dirs := []string{
		examplesDir,
		manDir + "/man1",
		manDir + "/man5",
	}

	for _, dir := range dirs {
		if err := ensureDirectory(dir); err != nil {
			return fmt.Errorf(
				"error ensuring the presence of %q: %w",
				dir,
				err,
			)
		}
	}

	if err := sh.Run(
		"go",
		"run",
		"./cmd/docgen",
		"--application-name="+appName(),
		"--binary-version="+binaryVersion(),
		"--man-dir="+manDir,
		"--examples-dir="+examplesDir,
	); err != nil {
		return fmt.Errorf("error generating the documentation: %w", err)
	}

	fmt.Println("Documentation successfully generated.")

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

// Codegen generates the CLI code from the CLI definitions file.
func Codegen() error {
	for _, pkg := range slices.All([]string{"executor", "cli"}) {
		internalPkg := "./internal/" + pkg

		fmt.Printf("Generating code for %q\n", internalPkg)

		if err := sh.Run("go", "generate", internalPkg); err != nil {
			return err
		}
	}

	fmt.Println("Code successfully generated.")

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
// If ENBAS_APP_COMMIT_REF is set, the value of that is returned, otherwise
// the latest git commit is returned.
func gitCommit() string {
	commit := os.Getenv(envAppCommitRef)
	if commit != "" {
		return commit
	}

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
	return filepath.Join(binDir, appName())
}

func appTitledName() string {
	runes := []rune(appName())
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// ensureDirectory checks to see if the specified directory is present.
// If it is not present then an attempt is made to create it.
func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("unable to create %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf(
				"received an unknown error after getting the directory information: %w",
				err,
			)
		}
	}

	return nil
}
