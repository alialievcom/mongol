package sites_core

import (
	"AliAlievMos/mongol/config"
	"AliAlievMos/mongol/mongo_connection"
	"AliAlievMos/mongol/router"
)

func StartApp(cfgPath string) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	collections := mongo_connection.ConnectMongo(cfg)
	router.RunRouter(collections, cfg)
}
