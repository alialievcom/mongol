package router

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/AliAlievMos/mongol/models"
	"github.com/AliAlievMos/mongol/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func createLoginHandler(collectionUsers *mongo.Collection, cfg *models.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			input  models.User
			dataDB models.User
		)

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"login": input.Login}
		result := collectionUsers.FindOne(c.Request.Context(), filter)
		if err := result.Err(); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		if err := result.Decode(&dataDB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user data"})
			return
		}

		hasher := sha1.New()
		hasher.Write([]byte(input.Password))
		hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(cfg.Api.SecretKey))

		if hashedPass != dataDB.Password {
			c.Status(http.StatusUnauthorized)
			return
		}

		custom := jwt.MapClaims{
			"username": input.Login,
			"exp":      time.Now().Add(500 * time.Hour).Unix(),
			"roles":    dataDB.Roles,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, custom)
		tokenString, err := token.SignedString(cfg.Api.SecretKey)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, models.Token{Access: tokenString, User: dataDB})
	}
}

func createCheckTokenHandler(cfg *models.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header["Token"]
		if len(token) == 0 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "no token in header")
			return
		}
		if token[0] == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "no token in header")
			return
		}
		roles, err := VerToken(
			c,
			token[0],
			cfg.Api.SecretKey,
		)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		// nil is equal all
		if roles == nil {
			return
		}
		method := c.Request.Method
		for _, r := range *roles {
			if method == r {
				return
			}
		}
		utils.ErrorResponse(c, http.StatusUnauthorized, "no such role")
	}
}

func VerToken(_ *gin.Context, tokenString string, secretKey []byte) (*[]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		unixTime, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("expired token")
		}
		timeExp := time.Unix(int64(unixTime), 0)
		validTime := time.Now().After(timeExp)
		if validTime {
			return nil, errors.New("expired token")
		}
		roles, ok := claims["roles"].([]string)
		if !ok {
			return nil, nil
		}
		return &roles, nil
	}
	return nil, errors.New("not valid token")
}

func insertUsers(userCollection *mongo.Collection, cfg *models.Config) {
	var existingUsers = make(map[string]struct{}, len(cfg.Mongo.Users))
	for _, user := range cfg.Mongo.Users {
		hasher := sha1.New()
		hasher.Write([]byte(user.Password))
		hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(cfg.Api.SecretKey))

		filter := bson.M{"login": user.Login}
		var existingUser models.User
		err := userCollection.FindOne(context.Background(), filter).Decode(&existingUser)

		if errors.Is(err, mongo.ErrNoDocuments) {
			user.ID = primitive.NewObjectID()
			user.Password = hashedPass

			_, err := userCollection.InsertOne(context.Background(), user)
			if err != nil {
				log.Panicf("Error inserting user %s: %v", user.Login, err)
			}
			existingUsers[user.Login] = struct{}{}
			continue
		} else if err != nil {
			log.Panicf("Error checking for user %s: %v", user.Login, err)
		}
		existingUser.Password = hashedPass
		existingUser.Roles = user.Roles
		update, err := utils.GenerateUpdateBson(existingUser)
		if err != nil {
			log.Panicf("Error GenerateUpdateBson for user %s: %v", user.Login, err)
		}
		_, err = userCollection.UpdateByID(context.Background(), existingUser.ID, update)
		if err != nil {
			log.Panicf("Error UpdateByID for user %s: %v", user.Login, err)
		}
		existingUsers[existingUser.Login] = struct{}{}
	}
}
