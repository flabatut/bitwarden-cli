/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/flabatut/bitwarden-cli/pkg/test"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("test called")
		w := &test.Workflow{
			Client:               daggerClient,
			ReleaseVersion:       viper.GetString("releaseVersion"),
			BuilderNodeJSVersion: viper.GetString("builderNodeJSVersion"),
			RunnerEntryPointPath: viper.GetString("runnerEntryPointPath"),
			BuilderWorkDir:       viper.GetString("builderWorkDir"),
			BuilderImage:         viper.GetString("builderImage"),
			RunnerImage:          viper.GetString("runnerImage"),
			BuilderPlatforms:     platforms,
		}
		if err := w.Test(cmd.Context()); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
