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

# Update New Version

Untuk membuat code yang lebih muda dibaca agar dapat dipergunakan lebih simple


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

	r.Run(":8080")
}

```


<h6> Contoh fungsi yang usang</h6>

```go
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
```


Jika ada kesalahan atau bug bisa menghubungi saya atau bikin issues pada repository ini
