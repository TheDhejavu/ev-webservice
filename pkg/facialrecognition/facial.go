package facialrecognition

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type facialRecognition struct {
	logger log.Logger
}

type FacialRecognition interface {
	Verify(id string, imageTmpPath string) (map[string]string, error) // Verify image
	Register(id string, imagesTmpPaths []string) error                // Register images
}

var (
	FACIAL_RECOGNITION_URL = "http://127.0.0.1:8000/api/auth"
)

func NewFacialRecogntion(logger log.Logger) FacialRecognition {
	return &facialRecognition{logger}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName string, paths []string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := 0; i < len(paths); i++ {
		path := paths[i]
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)

	}

	for key, val := range params {
		fmt.Println(params)
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
func (f facialRecognition) Verify(id string, imageTmpPath string) (res map[string]string, err error) {
	url := fmt.Sprintf("%s/verify", FACIAL_RECOGNITION_URL)
	extraParams := new(map[string]string)
	request, err := newfileUploadRequest(url, *extraParams, "image", []string{imageTmpPath})
	if err != nil {
		f.logger.Error(err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		f.logger.Error(err)
		return
	} else {
		body := &bytes.Buffer{}
		_, err = body.ReadFrom(resp.Body)
		if err != nil {
			f.logger.Error(err)
			return
		}
		resp.Body.Close()
		f.logger.Info("verify facial status code:", resp.StatusCode)
		f.logger.Info("verify facial header:", resp.Header)

		f.logger.Info("verify facial body", body)
		err = errors.New(body.String())
		if resp.StatusCode != 200 {
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return
		}
	}
	return
}
func (f facialRecognition) Register(id string, imageTmpPaths []string) error {
	url := fmt.Sprintf("%s/register", FACIAL_RECOGNITION_URL)

	extraParams := map[string]string{
		"user_id": id,
	}
	request, err := newfileUploadRequest(url, extraParams, "files", imageTmpPaths)
	if err != nil {
		f.logger.Error(err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		f.logger.Error(err)
		return err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			f.logger.Error(err)
			return err
		}
		resp.Body.Close()
		f.logger.Info("Register facial status code:", resp.StatusCode)
		f.logger.Info("Register facial header:", resp.Header)

		f.logger.Info("Register facial body", body)

		if resp.StatusCode != 200 {
			return errors.New(body.String())
		}
	}
	return nil
}
