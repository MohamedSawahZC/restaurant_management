package routes

import (
	"github.com/MohamedSawahZC/restaurant_management/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/invoices",controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id",controllers.GetInvoice())
	incomingRoutes.POST("/invoices",controllers.CreateInvoice())
	incomingRoutes.PATCH("/invoices/:invoice_id",controllers.UpdateInvoice())
}