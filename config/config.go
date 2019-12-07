package config

import (
	"go-dht/kademlia"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	ID            string
	Addr          string
	Port          string
	KnownNodesStr []KnownNodeStr        `mapstructure:"knownnodes"`
	KnownNodes    []kademlia.ListedNode `mapstructure:"-"`
}

type KnownNodeStr struct {
	ID   string
	Addr string
	Port string
}

var C Config

func MustRead(c *cli.Context) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&C); err != nil {
		return err
	}

	for _, v := range C.KnownNodesStr {
		id, err := kademlia.IDFromString(v.ID)
		if err != nil {
			return err
		}
		kn := kademlia.ListedNode{
			ID:   id,
			Addr: v.Addr,
			Port: v.Port,
		}
		C.KnownNodes = append(C.KnownNodes, kn)
	}
	C.KnownNodesStr = []KnownNodeStr{}
	return nil
}
