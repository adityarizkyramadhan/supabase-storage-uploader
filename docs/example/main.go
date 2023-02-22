package main

import (
	"os"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		link, err := supabasestorageuploader.Upload(os.Getenv("HOST"), os.Getenv("TOKEN"), os.Getenv("STORAGE_NAME"), os.Getenv("STORAGE_PATH"), file)
		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": link})
	})

	r.Run(":8080")
}
