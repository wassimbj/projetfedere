package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"pfserver/config"
	"pfserver/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func StartServer(Router *mux.Router, port string) error {

	handler := cors.New(config.Cors()).Handler(Router)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20, // = 1,048,576 bytes
	}

	return s.ListenAndServe()
}

// @params req.Body
// reads http request body, returns JSON
func ReadBody(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	s := string(data)
	return s, err
}

// respond to a request with a json
type ResOpts struct {
	Status int
	Msg    interface{}
}

func Respond(res http.ResponseWriter, opts ResOpts) (int, error) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(opts.Status)
	jsonData, _ := json.Marshal(opts.Msg)
	return res.Write(jsonData)
}

type FileDetails struct {
	OriginalName string // this is the original uploaded name
	SavedName    string // this is the saved file name
	// Size         int64
	Header textproto.MIMEHeader // MIME Header
}

type UploadFilesData struct {
	FileDetails []FileDetails
	Body        string // json
}

// fileKeys = input name (<input name="..." />) | map["input_name"] = file || text
// this function returns the uploaded files data, and also the req.Body json.
func UploadFiles(fileKeys map[string]bool, req *http.Request) (UploadFilesData, error) {
	mr, err := req.MultipartReader()
	if err != nil {
		return UploadFilesData{}, err
	}
	var data UploadFilesData
	timer := time.NewTimer(time.Second * 10) // timeout
	for {
		select {
		case <-timer.C:
			timer.Stop()
			return UploadFilesData{}, errors.New("Timed out !!")
		default:
			part, err := mr.NextPart()
			// This is OK, no more parts
			if err == io.EOF {
				return data, nil
			}

			if fileKeys[part.FormName()] {
				if part.FileName() == "" {
					// break, skip to the next loop
					continue
				}
				storageFolder := filepath.Join(utils.RootPath(), "storage")
				imgType := part.Header.Get("Content-Type")      // e.g: image/png
				fileExtension := strings.Split(imgType, "/")[1] // get file extension
				// * => will be replaced by the random string by the ioutil.TempFile function
				fileNameToUpload := fmt.Sprintf("%s-*.%s", uuid.NewString(), fileExtension)
				uploadedFilePath, err := utils.UploadTempFile(part, fileNameToUpload, storageFolder)
				// split by "/" and get the last elem which is basicly the filename /path/to/filename.ext
				uploadedFileName := strings.Split(uploadedFilePath, "/")[len(strings.Split(uploadedFilePath, "/"))-1]

				if err != nil {
					return UploadFilesData{}, err
				}
				data.FileDetails = append(data.FileDetails, FileDetails{
					OriginalName: part.FileName(),
					SavedName:    uploadedFileName,
					Header:       part.Header,
				})

			} else if part.FormName() == "body" {
				// body => json = {"doc_used": "..."}
				jsonBody, _ := ReadBody(part)
				data.Body = jsonBody
			}
		}
	}

}
