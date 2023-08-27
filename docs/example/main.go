package main

import (
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	// Buat Client
	supClient := supabasestorageuploader.New(
		"https://your-unique-url.supabase.co",
		"your-token",
		"your-bucket-name",
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

	r.GET("/list", func(c *gin.Context) {
		list, err := supClient.ListBucket(
			&supabasestorageuploader.RequestBodyListBucket{
				Limit:  10,
				Offset: 0,
				SortBy: struct {
					Column string `json:"column"`
					Order  string `json:"order"`
				}{
					Column: "name",
					Order:  "asc",
				},
			},
		)
		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": list})
	})

	r.DELETE("/delete", func(c *gin.Context) {
		// get body from request
		var requestBody map[string]string
		err := c.BindJSON(&requestBody)
		if err != nil {
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		err = supClient.Delete(requestBody["link"])
		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	})
	log.Printf("Server running at %v\n", color.GreenString("http://localhost:8080"))
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
