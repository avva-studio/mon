package cmd

import (
	"log"

	"github.com/glynternet/accounting-rest/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdDBServe = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := newStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error creating storage"))
		}
		s, err := server.New(store)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error creating new server"))
		}
		log.Fatal(s.ListenAndServe(":" + viper.GetString(keyPort)))
	},
}

func init() {
	rootCmd.AddCommand(cmdDBServe)
}
