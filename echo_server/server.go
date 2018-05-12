package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/*action", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": c.Param("action"),
		})
	})
	router.Run()
}
