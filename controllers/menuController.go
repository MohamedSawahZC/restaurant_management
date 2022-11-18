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
		err := menuCollection.FindOne(ctx,bson.M{"menu_id":menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK,menu)
		defer cancel()
	
	
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}

		validateError := validate.Struct(menu)

		if validateError != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validateError.Error()})
			return
		}
		defer cancel()
		menu.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		menu.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()
		result, inserErr := menuCollection.InsertOne(ctx,menu)
		if inserErr != nil {
			msg := fmt.Sprintf("menu item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}
		defer cancel()
		c.JSON(http.StatusCreated,result)
	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		menuId := c.Param("menu_id")
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var menu models.Menu 
		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		filter := bson.M{"menu_id":menuId}
		var updateObj primitive.D

		if menu.Start_Date !=nil && menu.End_Date !=nil{
			if !inTimeSpan(*menu.Start_Date,*menu.End_Date,time.Now()){
				msg := "Kindly retype the time please"
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
				defer cancel()
				return
			}
			updateObj = append(updateObj, bson.E{"start_date",menu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date",menu.End_Date})

			if menu.Name != ""{
				updateObj = append(updateObj, bson.E{"name",menu.Name})
			}
			if menu.Category != ""{
				updateObj = append(updateObj, bson.E{"category",menu.Category})
			}
			menu.Updated_at ,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_at",menu.Updated_at})
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			result,err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set",updateObj},
				},
				&opt,

			)
			if err != nil{
				msg :=	"Menu update failed"
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			}
			defer cancel()
			c.JSON(http.StatusOK,result)
		}


	}
}

