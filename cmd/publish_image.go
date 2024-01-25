/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/flabatut/bitwarden-cli/pkg/publish/image"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("publish image called")
		if err := runPublishImageCmd(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	publishCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPublishImageCmd(cmd *cobra.Command) error {
	containers, _, err := runBuildCmd(cmd)
	if err != nil {
		return err
	}
	username, err := getRegistryUsername()
	if err != nil {
		return err
	}
	password, err := getRegistryPassword()
	if err != nil {
		return err
	}

	job := &image.Workflow{
		Client:           daggerClient,
		ReleaseVersion:   viper.GetString("releaseVersion"),
		RegistryFQDN:     viper.GetString("registryFQDN"),
		ProjectNamespace: viper.GetString("projectNamespace"),
		RegistryUsername: username,
		RegistryPassword: password,
	}

	err = job.Publish(cmd.Context(), containers)
	if err != nil {
		return err
	}

	return err
}
