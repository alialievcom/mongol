package sites_core

import (
	"alialiev/sites-core/config"
	"alialiev/sites-core/mongo_connection"
	"alialiev/sites-core/router"
)

func StartApp(cfgPath string) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	collections := mongo_connection.ConnectMongo(cfg)
	router.RunRouter(collections, cfg)
}
