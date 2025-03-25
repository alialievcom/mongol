package router

import (
	"github.com/AliAlievMos/mongol/constants"
	sites_core "github.com/AliAlievMos/mongol/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func RunRouter(collections []*mongo.Collection, cfg *sites_core.Config) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Api.Origin},
		AllowMethods:     strings.Split(cfg.Api.Methods, ", "),
		AllowHeaders:     strings.Split(cfg.Api.Headers, ", "),
		AllowCredentials: true,
	}))
	r.Use(corsMiddleware(cfg.Api))
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
	authen := r.Group("/authenticated", createCheckTokenHandler(cfg))
	r.POST("/login", createLoginHandler(authCol, cfg))
	r.POST("/reg", createRegHandler(authCol, cfg))
	for _, collection := range collections {
		collectionName := collection.Name()
		details := cfg.GeneratedStructMap[collectionName]
		authen.POST("/"+collectionName, createCheckRoleMiddleware(details.NecessaryAuthRole[constants.Post]), createPostHandler(collection, details.Model))
		authen.PUT("/"+collectionName, createCheckRoleMiddleware(details.NecessaryAuthRole[constants.Put]), createPostHandler(collection, details.Model))
		r.GET("/"+collectionName, createCheckRoleMiddleware(details.NecessaryAuthRole[constants.Get]), createGetHandler(collection, details))
		authen.DELETE("/"+collectionName+"/:id", createCheckRoleMiddleware(details.NecessaryAuthRole[constants.Delete]), createDeleteHandler(collection))
	}

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	err := r.Run(cfg.Api.Port)
	if err != nil {
		panic(err)
	}
}
