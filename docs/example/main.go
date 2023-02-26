package main

import (
	"fmt"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Buat Client
	supClient := supabasestorageuploader.NewSupabaseClient(
		"PROJECT_URL",
		"PROJECT_API_KEYS",
		"PROJECT_STORAGE_NAME",
		"STORAGE_PATH",
	)

	r.POST("/upload/v2", func(c *gin.Context) {
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		link, err := supClient.Upload(file)
		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": link})
	})

	r.DELETE("file", func(c *gin.Context) {
		linkFile := c.Request.FormValue("linkfile")

		fmt.Println(linkFile)

		data, err := supClient.DeleteFile(linkFile)

		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": data})
	})

	r.Run(":8080")
}
