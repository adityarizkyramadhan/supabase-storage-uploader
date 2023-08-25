package supabasestorageuploader

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type Client struct {
	token      string
	httpClient *http.Client
	urlProject string
	fileUrl    string
	bucketName string
}

func New(
	projectUrl string,
	token string,
	bucketName string,
) *Client {
	url := projectUrl + "/storage/v1/object/" + bucketName + "/"
	fileUrl := projectUrl + "/storage/v1/object/public/" + bucketName + "/"
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		fileUrl:    fileUrl,
		urlProject: url,
		bucketName: bucketName,
	}
}

func (c *Client) Upload(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		log.Println("Error reading file:", errFileNotFound)
		return "", errFileNotFound
	}
	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Error opening file:", err)
		return "", err
	}
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		log.Println("Error creating form file:", err)
		return "", err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Println("Error copying file to form:", err)
		return "", err
	}
	err = multipartWriter.Close()
	if err != nil {
		log.Println("Error closing multipart writer:", err)
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, c.urlProject+fileHeader.Filename, &requestBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("Received non-200 response:", response.StatusCode)
		return "", err
	}
	return c.linkFile(fileHeader.Filename), nil
}

func (c *Client) linkFile(filename string) string {
	return c.fileUrl + filename
}
