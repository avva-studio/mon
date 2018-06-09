package cmd

import (
	"github.com/glynternet/mon/client"
	"github.com/spf13/viper"
)

func newClient() client.Client {
	return client.Client(viper.GetString(keyServerHost))
}
