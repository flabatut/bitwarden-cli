/*
Copyright © 2024 Franck Labatut
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	daggerClient *dagger.Client
	platforms    = []dagger.Platform{ // TODO: use viper with mapstruct
		"darwin/amd64", // a.k.a. x86_64
		"darwin/arm64", // a.k.a. aarch64
		"linux/amd64",  // a.k.a. x86_64
		"linux/arm64",  // a.k.a. aarch64
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitwarden-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initDaggerClient)
	cobra.OnFinalize(deferDaggerClient)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bitwarden-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bitwarden-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bitwarden-cli")
	}
	viper.SetEnvPrefix("BWCLI") // Support for env vars matching prefix below
	viper.AutomaticEnv()        // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initDaggerClient() {
	fmt.Println("Starting Dagger Engine session")
	var err error
	daggerClient, err = dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stderr))
	cobra.CheckErr(err)
}

func deferDaggerClient() {
	fmt.Println("Closing Dagger Engine session")
	err := daggerClient.Close()
	cobra.CheckErr(err)
}

func getRegistryUsername() (string, error) {
	// any viper field username found (works for env var BWCLI_USERNAME as well)
	if viper.IsSet("username") {
		return viper.GetString("username"), nil
	}
	// GITHUB_REPOSITORY_OWNER
	if username, ok := os.LookupEnv("GITHUB_REPOSITORY_OWNER"); ok {
		return username, nil
	}
	return "", fmt.Errorf("username for registry not found")
}

func getRegistryPassword() (*dagger.Secret, error) {
	// any viper field password found (works for env var BWCLI_PASSWORD as well)
	if viper.IsSet("password") {
		return daggerClient.SetSecret("password", viper.GetString("password")), nil
	}
	// support for GITHUB_TOKEN
	if password, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		return daggerClient.SetSecret("password", password), nil
	}
	return nil, fmt.Errorf("password for registry not found")
}
