package routes

import (
	"github.com/MohamedSawahZC/restaurant_management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orders",controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id",controllers.GetOrder())
	incomingRoutes.POST("/orders",controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:order_id",controllers.UpdateOrder())
}