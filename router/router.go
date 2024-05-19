package router

import (
	"alialiev/sites-core/constants"
	sites_core "alialiev/sites-core/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunRouter(collections []*mongo.Collection, cfg *sites_core.Config) {
	r := gin.Default()
	var userCol *mongo.Collection
	for _, collection := range collections {
		if collection.Name() == constants.UsersCollection {
			userCol = collection
		}
	}
	insertUsers(userCol, cfg)
	admin := r.Group("/admin", createCheckTokenHandler(cfg))
	r.POST("/login", createLoginHandler(userCol, cfg))
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
