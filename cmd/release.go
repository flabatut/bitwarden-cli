/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/flabatut/bitwarden-cli/pkg/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("publish artifact called")
		if _, err := runReleaseCmd(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	err := markReleaseVersionRequired(releaseCmd)
	cobra.CheckErr(err)
}

func runReleaseCmd(cmd *cobra.Command) (*dagger.Container, error) {
	fmt.Println("release called")
	password, err := getRegistryPassword()
	if err != nil {
		return nil, err
	}
	w := &release.Workflow{
		Client:           daggerClient,
		ProjectNamespace: viper.GetString("projectNamespace"),
		RegistryPassword: password,
	}
	if !viper.IsSet("releaseVersion") {
		return nil, fmt.Errorf("required flag(s) releaseVersion not set")
	}
	releaser, err := w.Release(cmd.Context(), viper.GetString("releaseVersion"))
	if err != nil {
		return nil, err
	}
	return releaser, nil
}
