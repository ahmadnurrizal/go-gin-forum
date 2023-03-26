package fileupload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, map[string]string)
}

// So what is exposed is Uploader
var FileUpload UploadFileInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, map[string]string) {
	errList := map[string]string{}

	// Validate file size
	size := file.Size
	if size > int64(1024000) {
		errList["Too_large"] = "Sorry, Please upload an Image of 1MB or less"
		return "", errList
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		errList["Not_Image"] = "Please upload a valid image (jpg, jpeg, png or gif)"
		return "", errList
	}

	// Open file for reading
	f, err := file.Open()
	if err != nil {
		errList["File_Open_Error"] = "Error opening file"
		return "", errList
	}
	defer f.Close()

	// Create local file for writing
	fileName := FormatFile(file.Filename)
	filePath := fileName

	localFile, err := os.Create("assets/upload/" + filePath)
	if err != nil {
		errList["Other_Err"] = "Failed to create file on server"
		return "", errList
	}
	defer localFile.Close()

	// Copy file contents to local file
	_, err = io.Copy(localFile, f)
	if err != nil {
		errList["Other_Err"] = "Failed to copy file contents to local file"
		return "", errList
	}

	// File uploaded successfully, return the local file path
	return fileName, nil
}
