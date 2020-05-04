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
)

// RecursiveZip is a function to create zip for given directory
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

// Backup function is to create cockroach database backup
func Backup(sqlFileNamePath string) error {
	dbName := os.Getenv("COCKROACH_DATABASE")
	log.Printf("%s database backup started.\n", dbName)

	dumpQueryString := fmt.Sprintf("/cockroach/cockroach dump %s ", dbName)

	connectionFlagsString := fmt.Sprintf("--user %s --host=%s", os.Getenv("COCKROACH_USER"), os.Getenv("COCKROACH_HOST"))

	if os.Getenv("COCKROACH_INSECURE") == "true" {
		connectionFlagsString = connectionFlagsString + " --insecure"
	}

	if os.Getenv("COCKROACH_CERTS_DIR") != "" {
		connectionFlagsString = connectionFlagsString + fmt.Sprintf(" --certs-dir=/cockroach-certs/")
	}

	query := dumpQueryString + connectionFlagsString + " > " + sqlFileNamePath

	// log.Println(query)

	cmd := exec.Command("/bin/sh", "-c", query)
	_, error := cmd.CombinedOutput()
	if error != nil {
		log.Fatalf("cmd.Run() failed with %s\n", error)
		return error
	}
	log.Printf("%s database backup done.\n", dbName)
	return nil
}
