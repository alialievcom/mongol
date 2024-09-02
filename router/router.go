package router

import (
	"github.com/AliAlievMos/mongol/constants"
	sites_core "github.com/AliAlievMos/mongol/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunRouter(collections []*mongo.Collection, cfg *sites_core.Config) {
	r := gin.Default()
	var authCol *mongo.Collection
	var authColName string
	if cfg.Mongo.Auth.AuthCollection != nil {
		authColName = *cfg.Mongo.Auth.AuthCollection
	} else {
		authColName = constants.UsersCollection
	}
	for _, collection := range collections {
		if collection.Name() == authColName {
			authCol = collection
			break
		}
	}
	if cfg.Mongo.Auth.AuthCollection == nil {
		insertUsers(authCol, cfg)
	}
	admin := r.Group("/admin", createCheckTokenHandler(cfg))
	r.POST("/login", createLoginHandler(authCol, cfg))
	for _, collection := range collections {
		collectionName := collection.Name()
		model := cfg.GeneratedStructMap[collectionName]

		admin.POST("/"+collectionName, createPostHandler(collection, model))
		r.GET("/"+collectionName, createGetHandler(collection))
		admin.DELETE("/"+collectionName+"/:id", createDeleteHandler(collection))
	}

	err := r.Run(cfg.Api.Port)
	if err != nil {
		panic(err)
	}
}
