/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/flabatut/bitwarden-cli/pkg/publish/artifact"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageCmd represents the image command
var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("publish artifact called")
		if err := runPublishArtifactCmd(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	publishCmd.AddCommand(artifactCmd)
	err := markReleaseVersionRequired(artifactCmd)
	cobra.CheckErr(err)
}

func runPublishArtifactCmd(cmd *cobra.Command) error {
	_, artifacts, err := runBuildCmd(cmd)
	if err != nil {
		return err
	}

	_, err = runReleaseCmd(cmd)
	if err != nil {
		return err
	}

	password, err := getRegistryPassword()
	if err != nil {
		return err
	}
	if !viper.IsSet("releaseVersion") {
		return fmt.Errorf("required flag(s) releaseVersion not set")
	}
	job := &artifact.Workflow{
		Client:            daggerClient,
		ReleaseVersion:    viper.GetString("releaseVersion"),
		RegistryFQDN:      viper.GetString("registryFQDN"),
		ProjectNamespace:  viper.GetString("projectNamespace"),
		RegistryPassword:  password,
		PlatformsVariants: platforms,
	}

	err = job.Publish(cmd.Context(), artifacts)
	if err != nil {
		return err
	}

	return err
}
