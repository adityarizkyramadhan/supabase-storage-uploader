# Supabase Storage Uploader

Tujuan untuk mengupload file ke supabase storage via golang dengan bantuan API dari javascript

# Cara Penggunaan

Download ekstensi
```
go get github.com/adityarizkyramadhan/supabase-storage-uploader
```

# Peraturan

- Maksimal file upload sebesar 3 * 1024 * 1024 byte
- Server API bersifat serverless sehingga harap maklum jika down atau lamban
- Jika merasa repo ini ada kekurangan bisa contact saya atau bikin issues
- Jika merasa repo ini berguna, bisa bantu star :) arigatouuu :)

# Update New Version

- v0.0.1 => Add upload file
- v0.0.2 => Untuk membuat code yang lebih muda dibaca agar dapat dipergunakan lebih simple
- v0.0.3 => Add delete file


```go
package main

import (
	"os"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Buat Client
	supClient := supabasestorageuploader.NewSupabaseClient(
		"PROJECT_URL",
		"PROJECT_API_KEYS",
		"STORAGE_NAME",
		"STORAGE_FOLDER",
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


	// Updates add delete file
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

```


Jika ada kesalahan atau bug bisa menghubungi saya atau bikin issues pada repository ini
