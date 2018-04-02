package main

import (
	"strings"

	"log"

	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyServerHost = "server-host"
	appName       = "accounting-cli"
)

var cmdRoot = &cobra.Command{
	Use: appName,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Run %s --help for usage\n", appName)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	cmdRoot.PersistentFlags().StringP(keyServerHost, "H", "", "server host")
	err := viper.BindPFlags(cmdRoot.PersistentFlags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding root command flags"))
	}
}

// initConfig sets AutomaticEnv in viper to true.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
