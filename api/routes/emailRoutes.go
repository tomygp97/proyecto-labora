package routes

import (
	"project/api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupEmailRoutes(r *gin.Engine) {
	emailController := controllers.NewEmailController()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Bienvenido a la API de correos electrónicos"})
	})

	r.GET("/emails", emailController.GetEmails)
	r.GET("/emails/:id", emailController.GetEmailByID)

	r.GET("/emails/search", emailController.SearchEmails)
}
