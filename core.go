package sites_core

import (
	"github.com/AliAlievMos/mongol/config"
	"github.com/AliAlievMos/mongol/mongo_connection"
	"github.com/AliAlievMos/mongol/router"
)

func StartApp(cfgPath string) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	collections := mongo_connection.ConnectMongo(cfg)
	router.RunRouter(collections, cfg)
}
