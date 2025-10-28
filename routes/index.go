package routes

import (
	"net/http"

	"github.com/catchnaren/go-scalable-servers/routes/handlers"
	"github.com/gin-gonic/gin"
)

func MountRoutes() *gin.Engine {
	handler := gin.Default()
	handler.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok from AIR",
		})
	})

	taskRoutes := handler.Group("/task")
	{
		taskRoutes.POST("/", handlers.SaveTask)
		taskRoutes.GET("/", handlers.ReadTask)
		taskRoutes.PATCH("/", handlers.UpdateTask)
		taskRoutes.DELETE("/:id", handlers.DeleteTask)
	}

	
	handler.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})
	return handler
}