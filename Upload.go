package main

import (
	"log"
	"os"

	"github.com/minio/minio-go"
)

// CleanUp is a function to delete uploaded files
func CleanUp(sourceFilePath string, destinationFilePath string) {
	os.Remove(sourceFilePath)
	os.Remove(destinationFilePath)

	log.Printf("%s file deleted after upload.\n", destinationFilePath)
}

// Upload is a function to upload zip to Spaces
func Upload(fileName string, sourceFilePath string, destinationFilePath string) {

	// check if file exist before upload and delete after uploading finished.
	if _, err := os.Stat(destinationFilePath); err == nil {
		s3Client, err := minio.New(os.Getenv("S3_URL"), os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), true)
		if err != nil {
			log.Fatalln(err)
		}
		if _, err := s3Client.FPutObject(os.Getenv("BUCKET_NAME"), fileName, destinationFilePath, minio.PutObjectOptions{
			ContentType: "application/zip",
		}); err != nil {
			log.Println(err)
		}
		log.Printf("%s uploaded successfully.\n", destinationFilePath)

		CleanUp(sourceFilePath, destinationFilePath)
	} else if os.IsNotExist(err) {
		log.Printf("%s file does not exist.\n", destinationFilePath)
	} else {
		log.Fatalln(err)
	}
}
