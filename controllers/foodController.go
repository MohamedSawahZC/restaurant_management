package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MohamedSawahZC/restaurant_management/database"
	"github.com/MohamedSawahZC/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client,"food")

var validate = validator.New()
func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}
func GetFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food
		err := foodCollection.FindOne(ctx,bson.M{"food_id":foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK,food)
		defer cancel()
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var food models.Food
		var menu models.Menu

		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}

		validateError := validate.Struct(food)

		if validateError != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validateError.Error()})
			return
		}
		err := menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil{
			msg := fmt.Sprintf("Menu was not found")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		food.Created_at,_ = time.Parse(time.RFC3339,time.Now()).Format(time.RFC3339)
		food.Updated_at,_ = time.Parse(time.RFC3339,time.Now()).Format(time.RFC3339)
		food.ID = primitive.NewObjectID()
		food.Food_id=food.ID.Hex()
		var num = toFixed(*food.Price,2)
		food.Price = &num
		result, inserErr := foodCollection.InsertOne(ctx,food)
		if inserErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}
		defer cancel()
		c.JSON(http.StatusCreated,result)

	}
}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}


func round(num float64) int{

}

func toFixed(num float64,precision int) float64{

}