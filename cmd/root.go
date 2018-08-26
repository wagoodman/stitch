// Copyright © 2018 Alex Goodman
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wagoodman/stitch/core"
)

var cfgFile, cfgDir, homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stitch",
	Short: "A developer tool for making multi-repo project orchestration a breeze",
	Long:  `A developer tool for making multi-repo project orchestration a breeze.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// find home directory
	var err error
	homeDir, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfgDir = filepath.Join(homeDir, ".stitch")

	cobra.OnInitialize(initConfig, initState)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stitch.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// search config in ~/.stitch/ directory with name "stitch" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("stitch")
	}

	// read in environment variables that match
	viper.SetEnvPrefix("stitch")
	viper.AutomaticEnv()

	// set any defaults
	viper.SetDefault("workspace-path", filepath.Join(homeDir, "stitch"))

	// if a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// todo: check if verbose is set
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// if the config does not exist, create it
		finalCfgPath := filepath.Join(cfgDir, "stitch.yaml")
		if err := viper.SafeWriteConfigAs(finalCfgPath); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(finalCfgPath)
			}
		}
	}

	// override any values you don't want controlled from the flat file
	viper.Set("state-path", filepath.Join(cfgDir, "workspace.gob"))
}

// initState creates the directory for all application state to be persisted (based on current user account)
func initState() {
	// create config and workspace dir's if they don't exist
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		os.Mkdir(cfgDir, 0755)
	}

	workspaceDir := viper.GetString("workspace-path")
	if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
		os.Mkdir(workspaceDir, 0755)
	}

	// load the initial config from disk
	core.GetWorkspace()
}
