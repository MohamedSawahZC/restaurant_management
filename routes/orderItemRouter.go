package routes

import (
	"github.com/MohamedSawahZC/restaurant_management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderItems",controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id",controllers.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id",controllers.GetOrderItems())
	incomingRoutes.POST("/orderItems",controllers.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItem_id",controllers.UpdateOrderItem())
}