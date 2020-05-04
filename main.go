package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron"
)

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
)

// PostHandler converts post request body to string
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		log.Println("Manual backup requested.")
		AppFunction()
		log.Println("Manual backup completed.")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

// AppFunction is a function to run the application
func AppFunction() {
	fileName := fmt.Sprintf("%s_backup_%d", os.Getenv("COCKROACH_DATABASE"), time.Now().Unix())
	sqlFileNamePath := "/data/" + fileName + ".sql"

	Backup(sqlFileNamePath)

	zipFileName := fileName + ".tar.gz"

	zipDestinationFilePath := "/data/" + zipFileName
	RecursiveZip(sqlFileNamePath, zipDestinationFilePath)

	log.Printf("%s zip created.\n", zipDestinationFilePath)

	if os.Getenv("ACCESS_KEY_ID") != "" {
		if os.Getenv("BUCKET_NAME") == "" {
			log.Fatalln("BUCKET_NAME can't be blank.")
		}
		if os.Getenv("S3_URL") == "" {
			log.Fatalln("S3_URL can't be blank")
		}
		if os.Getenv("SECRET_ACCESS_KEY") == "" {
			log.Fatalln("SECRET_ACCESS_KEY can't be blank.")
		}
		Upload(zipFileName, sqlFileNamePath, zipDestinationFilePath)
	}
}

func main() {
	if os.Getenv("COCKROACH_DATABASE") == "" {
		log.Fatalln("Database env is not defined. Please set is by passing -e COCKROACH_DATABASE=.")
	}
	if os.Getenv("CRON_SCHEDULE") == "" {
		log.Fatalln("CRON_SCHEDULE can't be blank.")
	}
	if os.Getenv("COCKROACH_USER") == "" {
		os.Setenv("COCKROACH_USER", "root")
	}
	c := cron.New()
	c.AddFunc(os.Getenv("CRON_SCHEDULE"), AppFunction)
	c.Start()
	log.Println("Backup scheduler started successfully.")

	mux := http.NewServeMux()
	mux.HandleFunc("/backup", PostHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}
