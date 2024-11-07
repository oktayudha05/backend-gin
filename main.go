package main

import (
	"github.com/gin-gonic/gin"

	"github.com/oktayudha05/backend-gin/controllers"
)

func main(){
	router := gin.Default()

	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	mahasiswaGroup := router.Group("/mahasiswa")
	{
		mahasiswaGroup.GET("/", controllers.GetMahasiswa)
		mahasiswaGroup.GET("/:npm", controllers.GetMahasiswaByNpm)
		mahasiswaGroup.POST("/", controllers.PostMahasiswa)
		mahasiswaGroup.DELETE("/:npm", controllers.DeleteMahasiswaByNpm)
	}

	router.Run(":3000")
}
