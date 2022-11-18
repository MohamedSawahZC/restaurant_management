package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/MohamedSawahZC/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}
func GetTable() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		tableId := c.Param("menu_id")
		var table models.Table
		err := foodCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK,table)
		defer cancel()
	
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

