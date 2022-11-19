package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MohamedSawahZC/restaurant_management/database"
	"github.com/MohamedSawahZC/restaurant_management/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")


func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		result , err := orderCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the order"})
		}
		var allOrders []bson.M

		if err = result.All(ctx,&allOrders); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allOrders)

	}
}
func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order
		err := menuCollection.FindOne(ctx,bson.M{"order_id":orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK,order)
		defer cancel()
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var order models.Order
		var table models.Table
		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}

		validateError := validate.Struct(order)

		if validateError != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validateError.Error()})
			return
		}
		if order.Table_id !=nil{
			err := tableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message: table was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				return
			}
		}
		defer cancel()
		order.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		result, inserErr := orderCollection.InsertOne(ctx,order)
		if inserErr != nil {
			msg := fmt.Sprintf("menu item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}
		defer cancel()
		c.JSON(http.StatusCreated,result)
	}
	
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		orderId := c.Param("order_id")
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var order models.Order 
		var table models.Table
		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		if order.Table_id !=nil{
			err := tableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message: table was not found")
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			}
		}
		var updateObj primitive.D
		updateObj = append(updateObj,bson.E{"table_id",order.Table_id})
		order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj =  append(updateObj,bson.E{"updated_at",order.Updated_at})
		upsert := true
		filter := bson.M{
			"order_id":orderId,
		}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result,err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,result)

		}
}

