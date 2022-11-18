package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MohamedSawahZC/restaurant_management/database"
	"github.com/MohamedSawahZC/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		results , err := menuCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error occured while listing menus"})
		}
		var allMenus []bson.M
		if err = results.All(ctx,&allMenus); err != nil{
			log.Fatal(err)
		}
	}
}
func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := foodCollection.FindOne(ctx,bson.M{"menu_id":menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK,menu)
		defer cancel()
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

