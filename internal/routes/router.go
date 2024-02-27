package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lits-06/manage-user/internal/routes/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout",  handlers.Logout)
	r.GET("/getrecords/:email", handlers.Showinfo)
}