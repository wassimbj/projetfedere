package utils

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"
)

// validate the fields
var validate *validator.Validate

func Validator() *validator.Validate {

	if validate == nil {
		validate = validator.New()
	}

	return validate
}

func RootPath() string {
	_, file, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(file)

	return filepath.Join(basepath, "..") // get out of the /utils dir, which get us the root

}

func UploadTempFile(file io.Reader, tempFilePattren string, uploadPath string) (string, error) {
	// tempFilePattren := uuid.NewString() + "*.png"
	tempFile, err := ioutil.TempFile(uploadPath, tempFilePattren)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	return tempFile.Name(), nil

}
