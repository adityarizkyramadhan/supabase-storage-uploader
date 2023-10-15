package supabasestorageuploader

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
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
	fileUrl := projectUrl + "/storage/v1/object/public/" + bucketName + "/"
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		fileUrl:    fileUrl,
		urlProject: projectUrl,
		bucketName: bucketName,
	}
}

func (c *Client) Upload(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		log.Printf("%v %v \n", color.RedString("Error reading file header:"), ErrFileNotFound)
		return "", ErrFileNotFound
	}
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error opening file:"), err)
		return "", err
	}
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error creating form file:"), err)
		return "", err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error copying file:"), err)
		return "", err
	}
	err = multipartWriter.Close()
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error closing multipart writer:"), err)
		return "", err
	}
	url := c.urlProject + "/storage/v1/object/" + c.bucketName + "/" + fileHeader.Filename
	request, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error creating request:"), err)
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error sending request:"), err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v %v \n", color.RedString("Error closing response body:"), err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Printf("%v %v \n", color.RedString("Received non-200 response:"), response.StatusCode)
		return "", ErrBadRequest
	}
	link := c.linkFile(fileHeader.Filename)
	reqImage, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error creating request:"), err)
		return "", err
	}
	resImage, err := c.httpClient.Do(reqImage)
	if err != nil {
		return "", err
	}
	if resImage.StatusCode != http.StatusOK {
		log.Printf("%v %v \n", color.RedString("Received non-200 response:"), resImage.StatusCode)
		return "", ErrBadRequest
	}
	return link, nil
}

func (c *Client) linkFile(filename string) string {
	return c.fileUrl + filename
}

type RequestBodyListBucket struct {
	Prefix string `json:"prefix"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy struct {
		Column string `json:"column"`
		Order  string `json:"order"`
	} `json:"sortBy"`
	Search string `json:"search"`
}

type ResponseListBucket []struct {
	Name           string    `json:"name"`
	ID             string    `json:"id"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	Metadata       struct {
		CacheControl   string    `json:"cacheControl"`
		ContentLength  int       `json:"contentLength"`
		ETag           string    `json:"eTag"`
		HTTPStatusCode int       `json:"httpStatusCode"`
		LastModified   time.Time `json:"lastModified"`
		Mimetype       string    `json:"mimetype"`
		Size           int       `json:"size"`
	} `json:"metadata"`
}

func (c *Client) ListBucket(requestBody *RequestBodyListBucket) (*ResponseListBucket, error) {
	if requestBody == nil {
		log.Printf("%v %v \n", color.RedString("Error reading request body:"), ErrFileNotFound)
		return nil, ErrFileNotFound
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error marshalling request body:"), err)
		return nil, err
	}
	url := c.urlProject + "/storage/v1/object/list/" + c.bucketName
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error creating request:"), err)
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error sending request:"), err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v %v \n", color.RedString("Error closing response body:"), err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Printf("%v %v \n", color.RedString("Received non-200 response:"), response.StatusCode)
		return nil, err
	}
	var responseBody ResponseListBucket
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error decoding response body:"), err)
		return nil, err
	}
	return &responseBody, nil
}

func (c *Client) extractFilename(link string) string {
	return strings.ReplaceAll(link, c.fileUrl, "")
}

func (c *Client) Delete(link string) error {
	fileName := c.extractFilename(link)
	url := c.urlProject + "/storage/v1/object/" + c.bucketName + "/" + fileName
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("%v %v \n", color.RedString("Error creating request:"), err)
		return err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Fatalf("%v %v \n", color.RedString("Error sending request:"), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v %v \n", color.RedString("Error closing response body:"), err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Printf("%v %v \n", color.RedString("Received non-200 response:"), response.StatusCode)
		return ErrBadRequest
	}
	return nil
}
