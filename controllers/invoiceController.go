package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/MohamedSawahZC/restaurant_management/database"
	"github.com/MohamedSawahZC/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct{
	Invoice_id	string		
	Order_id	string 				
	Payment_method string		
	Payment_status *string		
	Payment_due_date interface{}
	Table_number interface{}
	Order_details interface{}
}
var invoiceColection *mongo.Collection = database.OpenCollection(database.Client,"invoice")
func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		results , err := invoiceColection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error occured while listing invoices"})
		}
		var allInvoices []bson.M
		if err = results.All(ctx,&allInvoices); err != nil{
			log.Fatal(err)
		}
	}
}
func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceColection.FindOne(ctx,bson.M{"invoice_id":invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the menu item"})
		}
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = "null"
		if invoice.Payment_method !=nil{
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoiceView.Payment_status
		invoiceView.Payment_due_date = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]
		defer cancel()
		c.JSON(http.StatusCreated,invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var invoice models.Invoice 
		if err := c.BindJSON(&invoice); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		validateError := validate.Struct(invoice)

		if validateError != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validateError.Error()})
			return
		}
		var order models.Order
		err := orderCollection.FindOne(ctx,bson.M{"order_id":invoice.Order_id}).Decode(&order)
		if err != nil{
			msg := fmt.Sprintf("Error while finding order")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		status := "PENDING"
		if invoice.Payment_status ==nil{
			invoice.Payment_status = &status
		}
		invoice.Payment_due_date,_ = time.Parse(time.RFC3339,time.Now().AddDate(0,0,1).Format(time.RFC3339))
		invoice.Created_at ,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		invoice.Updated_at ,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()
		result, inserErr := menuCollection.InsertOne(ctx,order)
		if inserErr != nil {
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
		}
		defer cancel()
		c.JSON(http.StatusCreated,result)
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		invoiceId := c.Param("invoice_id")
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		var invoice models.Invoice 
		if err := c.BindJSON(&invoice); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		filter := bson.M{"invoice_id": invoiceId}
		var updateObj primitive.D
		if invoice.Payment_method !=nil{
			updateObj=append(updateObj, bson.E{"payment_method",invoice.Payment_method})
		}
		if invoice.Payment_status !=nil{
			updateObj = append(updateObj, bson.E{"payment_status",invoice.Payment_status})
		}
		invoice.Updated_at ,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at",invoice.Updated_at})
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.Payment_status ==nil{
			invoice.Payment_status = &status
		}
		results,err := invoiceColection.UpdateOne(
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
		c.JSON(http.StatusOK,results)

	}
}

