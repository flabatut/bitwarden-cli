/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		w := &build.Workflow{
			Client:               daggerClient,
			PublishAddr:          viper.GetString("publishAddr"),
			ReleaseVersion:       viper.GetString("releaseVersion"),
			BuilderNodeJSVersion: viper.GetString("builderNodeJSVersion"),
			RunnerEntryPointPath: viper.GetString("runnerEntryPointPath"),
			BuilderWorkDir:       viper.GetString("builderWorkDir"),
			BuilderImage:         viper.GetString("builderImage"),
			RunnerImage:          viper.GetString("runnerImage"),
			BuilderPlatforms:     platforms,
		}
		if err := w.Build(cmd.Context()); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.SetDefault("builderWorkDir", "/build")                                                      // same as Dockerfile WORKDIR
	viper.SetDefault("builderImage", "mcr.microsoft.com/devcontainers/typescript-node:1-20-bullseye") // same as Dockerfile FROM, image for builder container
	viper.SetDefault("builderNodeJSVersion", "latest")                                                // vercel compatible format https://github.com/vercel/pkg
	viper.SetDefault("runnerImage", "docker.io/debian:bullseye-slim")                                 // same as Dockerfile FROM, image for final target container
	viper.SetDefault("runnerEntryPointPath", "/entrypoint")                                           // same as Dockerfile ENTRYPOINT
	viper.SetDefault("releaseVersion", "v2024.1.0")
	viper.SetDefault("publishAddr", "ghcr.io/flabatut/bitwarden-cli:latest")
	// if using local registry (https://docs.dagger.io/252029/load-images-local-docker-engine/#approach-2-use-a-local-registry-server)
	// viper.SetDefault("publishAddr", "localhost:5000/bitwarden-cli:latest")
}
