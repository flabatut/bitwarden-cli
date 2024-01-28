/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/flabatut/bitwarden-cli/pkg/build"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and push docker images and binaries",
	Long: `- build binaries (all os/arch)
	- build docker images (all arch)
	- push docker images to ghcr registry
	- push binaries to github`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("build called")
		if _, _, err := runBuildCmd(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	viper.SetDefault("builderWorkDir", "/build") // same as Dockerfile WORKDIR
	// viper.SetDefault("builderImage", "mcr.microsoft.com/devcontainers/typescript-node:1-20-bullseye") // same as Dockerfile FROM, image for builder container
	viper.SetDefault("builderImage", "mcr.microsoft.com/devcontainers/typescript-node:18-bullseye") // same as Dockerfile FROM, image for builder container

	viper.SetDefault("builderNodeJSVersion", "latest")                // vercel compatible format https://github.com/vercel/pkg
	viper.SetDefault("runnerImage", "docker.io/debian:bullseye-slim") // same as Dockerfile FROM, image for final target container
	viper.SetDefault("runnerEntryPointPath", "/entrypoint")           // same as Dockerfile ENTRYPOINT
	viper.SetDefault("registryFQDN", "ghcr.io")                       // TODO: revamp, remove publishaddress
	viper.SetDefault("projectNamespace", "flabatut/bitwarden-cli")    // TODO: make sure no / at the begining , discover value
	// viper.SetDefault("releaseVersion", "v2024.1.0")

	err := markReleaseVersionRequired(buildCmd)
	cobra.CheckErr(err)
}

func runBuildCmd(cmd *cobra.Command) ([]*dagger.Container, *dagger.Directory, error) {
	if !viper.IsSet("releaseVersion") {
		return nil, nil, fmt.Errorf("required flag(s) releaseVersion not set")
	}
	job := &build.Workflow{
		Client:               daggerClient,
		ReleaseVersion:       viper.GetString("releaseVersion"),
		BuilderNodeJSVersion: viper.GetString("builderNodeJSVersion"),
		RunnerEntryPointPath: viper.GetString("runnerEntryPointPath"),
		BuilderWorkDir:       viper.GetString("builderWorkDir"),
		BuilderImage:         viper.GetString("builderImage"),
		RunnerImage:          viper.GetString("runnerImage"),
		BuilderPlatforms:     platforms,
	}

	containers, artifacts, err := job.Build(cmd.Context())
	if err != nil {
		return nil, nil, err
	}
	return containers, artifacts, nil
}
