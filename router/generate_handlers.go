package router

import (
	"github.com/AliAlievMos/mongol/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"reflect"
)

func createPostHandler(collection *mongo.Collection, model reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		pub := reflect.New(model).Interface()

		if err := c.ShouldBindJSON(pub); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if pub == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating update BSON"})
			return
		}
		pubValue := reflect.ValueOf(pub).Elem()
		idField := pubValue.FieldByName("ID")

		if !idField.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model: missing ID field"})
			return
		}

		if idField.IsNil() {
			oid := primitive.NewObjectID()
			idField.Set(reflect.ValueOf(&oid))
		}

		update, err := utils.GenerateUpdateBson(pubValue.Interface())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating update BSON"})
			return
		}

		opt := options.Update().SetUpsert(true)
		result, err := collection.UpdateByID(c.Request.Context(), idField.Interface(), update, opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating or creating document"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func createGetHandler(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var files []bson.M
		id := c.Query("id")

		var filter bson.M
		if id != "" {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
				return
			}
			filter = bson.M{"_id": objID}
		} else {
			filter = bson.M{}
		}

		cursor, err := collection.Find(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching documents"})
			return
		}
		defer cursor.Close(c.Request.Context())

		if err = cursor.All(c.Request.Context(), &files); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading documents"})
			return
		}

		if files == nil {
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

func createDeleteHandler(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		_, err = collection.DeleteOne(c.Request.Context(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting document"})
			return
		}

		c.Status(http.StatusOK)
	}
}
