package supabasestorageuploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func Upload(host, token, storageName, storagePath string, fileHeader *multipart.FileHeader) (string, error) {

	url := "https://express-uploader-two.vercel.app/upload"
	method := "POST"

	file, err := fileHeader.Open()
	if err != nil {
		panic(err.Error())
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("host", host)
	_ = writer.WriteField("token", token)
	_ = writer.WriteField("storage_name", storageName)
	_ = writer.WriteField("storage_path", storagePath)

	defer file.Close()
	part5, err := writer.CreateFormFile("avatar", fileHeader.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(part5, file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	response := make(map[string]interface{})

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return response["data"].(string), nil
}
