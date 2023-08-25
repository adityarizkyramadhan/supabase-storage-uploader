package supabasestorageuploader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

var (
	errFileNotFound     = errors.New("fileHeader is null")
	errFileNotInStorage = errors.New("file not found, check your storage name, file path, and file name")
	errLinkNotFound     = errors.New("file not found, check your storage name, file path, file name, and policy")
)

const (
	urlPost   = "https://express-uploader-two.vercel.app/upload"
	urlDelete = "https://express-uploader-two.vercel.app/file"
	post      = "POST"
	delete    = "DELETE"
)

type (
	supabaseClient struct {
		ProjectUrl         string
		ProjectApiKeys     string
		ProjectStorageName string
		StorageFolder      string
	}
	SupabaseClientService interface {
		Upload(fileHeader *multipart.FileHeader) (string, error)
		DeleteFile(link string) (interface{}, error)
	}
	payload struct {
		Host        string `json:"host"`
		Token       string `json:"token"`
		StorageName string `json:"storage_name"`
		Link        string `json:"link"`
	}
)

func NewSupabaseClient(
	projectUrl string,
	projectApiKeys string,
	projectStorageName string,
	storageFolder string,
) SupabaseClientService {
	return &supabaseClient{
		ProjectUrl:         projectUrl,
		ProjectApiKeys:     projectApiKeys,
		ProjectStorageName: projectStorageName,
		StorageFolder:      storageFolder,
	}
}

func (sc *supabaseClient) Upload(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", errFileNotFound
	}

	file, err := fileHeader.Open()
	if err != nil {
		panic(err.Error())
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("host", sc.ProjectUrl)
	_ = writer.WriteField("token", sc.ProjectApiKeys)
	_ = writer.WriteField("storage_name", sc.ProjectStorageName)
	_ = writer.WriteField("storage_path", sc.StorageFolder)

	defer file.Close()
	part, err := writer.CreateFormFile("avatar", fileHeader.Filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest(post, urlPost, payload)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response := make(map[string]interface{})

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	err = sc.checkLink(response["data"].(string))
	if err != nil {
		return "", err
	}

	return response["data"].(string), nil

}

func (sc *supabaseClient) checkLink(link string) error {
	// Create an HTTP client
	client := &http.Client{}
	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return err
	}

	// Perform the request via the client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return errLinkNotFound
	}

	return nil
}

func (sc *supabaseClient) DeleteFile(link string) (interface{}, error) {
	payloadRequest := payload{
		Host:        sc.ProjectUrl,
		StorageName: sc.ProjectStorageName,
		Token:       sc.ProjectApiKeys,
		Link:        link,
	}
	requestBody, err := json.Marshal(&payloadRequest)
	if err != nil {
		fmt.Println("marshal json")
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(delete, urlDelete, bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(string(body))
	response := make(map[string]interface{})

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response["data"].([]interface{})) == 0 {
		return nil, errFileNotInStorage
	}

	return response["data"], nil
}
