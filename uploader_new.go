package supabasestorageuploader

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

var (
	errFileNotFound = errors.New("fileHeader is null")
)

const (
	url    = "https://express-uploader-two.vercel.app/upload"
	method = "POST"
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
	req, err := http.NewRequest(method, url, payload)

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

	return response["data"].(string), nil

}
