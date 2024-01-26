/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	// "github.com/flabatut/bitwarden-cli/pkg/publish/image"
	// "github.com/flabatut/bitwarden-cli/pkg/publish/artifact"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// imageCmd represents the image command
var publishAllCmd = &cobra.Command{
	Use:   "all",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("publish all called")
		if err := runPublishAllCmd(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	publishCmd.AddCommand(publishAllCmd)
	err := markReleaseVersionRequired(publishAllCmd)
	cobra.CheckErr(err)
}

func runPublishAllCmd(cmd *cobra.Command) error {
	if err := runPublishArtifactCmd(cmd); err != nil {
		return err
	}
	if err := runPublishImageCmd(cmd); err != nil {
		return err
	}
	return nil
}
