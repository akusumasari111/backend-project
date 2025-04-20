package routes

import (
	"backend-project/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/products/:id/upload", handlers.UploadProductImage)
	r.GET("/products/:id/image", handlers.DownloadProductImage)
	r.DELETE("/products/:id/image", handlers.DeleteProductImage)
}
