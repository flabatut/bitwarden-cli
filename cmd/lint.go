/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/flabatut/bitwarden-cli/pkg/lint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("lint called")

		w := &lint.Workflow{
			Client:        daggerClient,
			LinterWorkDir: viper.GetString("linterWorkDir"),
			LinterImage:   viper.GetString("linterImage"),
		}
		if err := w.Lint(cmd.Context()); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetDefault("LinterWorkDir", "/build")                                 // same as Dockerfile WORKDIR
	viper.SetDefault("LinterImage", "docker.io/golangci/golangci-lint:v1.55.2") // same as Dockerfile FROM, image for builder container

}
