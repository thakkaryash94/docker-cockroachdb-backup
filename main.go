package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go"
	"github.com/robfig/cron"
)

// Backup function to create cockroach database backup
func Backup() error {
	log.Printf("%s database backup started.\n", os.Getenv("COCKROACH_DATABASE"))
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("/cockroach/cockroach dump %s --insecure --host=%s > /data/backup.sql", os.Getenv("COCKROACH_DATABASE"), os.Getenv("COCKROACH_HOST")))
	_, error := cmd.CombinedOutput()
	if error != nil {
		log.Fatalf("cmd.Run() failed with %s\n", error)
		return error
	}
	log.Printf("%s database backup done.\n", os.Getenv("COCKROACH_DATABASE"))
	return nil
}

// RecursiveZip is function to create zip for given directory
func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

// Upload is function to upload zip to Spaces
func Upload() {

	const sourceFilePath = "/data/backup.sql"

	Backup()

	fileName := fmt.Sprintf("backup_%d.tar.gz", time.Now().Unix())
	destinationFilePath := "/data/" + fileName
	RecursiveZip(sourceFilePath, destinationFilePath)

	log.Printf("%s zip created.\n", destinationFilePath)

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

		os.Remove(sourceFilePath)
		os.Remove(destinationFilePath)

		log.Printf("%s file deleted after upload.\n", destinationFilePath)

	} else if os.IsNotExist(err) {
		log.Printf("%s file does not exist.\n", destinationFilePath)
	} else {
		log.Fatalln(err)
	}
}

func main() {
	if os.Getenv("ACCESS_KEY_ID") == "" {
		log.Fatalln("ACCESS_KEY_ID can't be blank.")
	}
	if os.Getenv("BUCKET_NAME") == "" {
		log.Fatalln("BUCKET_NAME can't be blank.")
	}
	if os.Getenv("CRON_SCHEDULE") == "" {
		log.Fatalln("CRON_SCHEDULE can't be blank.")
	}
	if os.Getenv("S3_URL") == "" {
		log.Fatalln("S3_URL can't be blank")
	}
	if os.Getenv("SECRET_ACCESS_KEY") == "" {
		log.Fatalln("SECRET_ACCESS_KEY can't be blank.")
	}
	c := cron.New()
	c.AddFunc(os.Getenv("CRON_SCHEDULE"), Upload)
	c.Start()
	log.Println("Backup scheduler started successfully.")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
