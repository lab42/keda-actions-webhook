//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
)

// Project setup
func Init() error {
	fmt.Println("Configuring git hooks")
	if err := sh.RunV("git", "config", "core.hooksPath", ".githooks"); err != nil {
		return err
	}
	return nil
}

// Dependency installation
func Install() error {
	fmt.Println("Installing svu")
	if err := sh.RunV("go", "install", "github.com/caarlos0/svu@latest"); err != nil {
		return err
	}

	fmt.Println("Installing staticcheck")
	if err := sh.RunV("go", "install", "honnef.co/go/tools/cmd/staticcheck@latest"); err != nil {
		return err
	}

	fmt.Println("Installing golangci-lint")
	if err := sh.RunV("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1"); err != nil {
		return err
	}

	fmt.Println("Installing air")
	if err := sh.RunV("go", "install", "github.com/cosmtrek/air@latest"); err != nil {
		return err
	}

	return nil
}

// Start development with hot reload
func Dev() error {
	if err := sh.RunV("air"); err != nil {
		return err
	}
	return nil
}

// Run configuration check, fmt and linting
func Validate() error {
	fmt.Println("Formatting code")
	if err := sh.RunV("go", "fmt", "./..."); err != nil {
		return err
	}

	fmt.Println("Checking go modules")
	if err := sh.RunV("go", "mod", "verify"); err != nil {
		return err
	}

	fmt.Println("Checking code")
	if err := sh.RunV("staticcheck", "./..."); err != nil {
		return err
	}

	fmt.Println("Linting code")
	if err := sh.RunV("golangci-lint", "run", "./..."); err != nil {
		return err
	}

	fmt.Println("Checking goreleaser configuration")
	if err := sh.RunV("goreleaser", "check"); err != nil {
		return err
	}
	return nil
}

// Builds the current project
func Build() error {
	if err := sh.RunV("goreleaser", "release", "--snapshot", "--rm-dist", "--skip-publish"); err != nil {
		return err
	}
	return nil
}
