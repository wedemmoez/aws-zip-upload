//Upload files in containing folder to IR upload bucket

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var uploader *s3manager.Uploader
var filename string

// create uploader to load ZIP file to S3
func newUploader() *s3manager.Uploader {

	s3Config := &aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	}
	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)
	return uploader
}

// upload ZIP to AWS.
func upload() {
	// Test upload using "test.zip" archive containing two XKCD images
	// filename = "test.zip"
	log.Println("uploading " + filename)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3_BUCKET), // bucket's name
		Key:    aws.String(filename),  // files destination location
		Body:   f,                     // file content
	})

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
}

func main() {
	//get hostname of device
	hostname, err := (os.Hostname())
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
	}
	//create filename using hostname and current datetime in specific format
	curTime := time.Now().Format("2006-01-02_1504")
	filename = hostname + "_" + curTime + ".zip"
	fmt.Println("creating zip archive with name " + filename)

	//create ZIP archive
	archive, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	w := zip.NewWriter(archive)
	defer w.Close()

	//create os.walker function to loop through directory and put all files into ZIP
	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		// fmt.Println(info.Name())
		if err != nil {
			return err
		}
		//ignore the zip file we're adding into, otherwise it'll try to add the existing zip to itself over and over
		if info.Name() == filename {
			return nil
		}

		//ignore IR Upload executable
		if info.Name() == "awsupload.exe" {
			return nil
		}

		if info.IsDir() {
			// add a trailing slash for creating dir
			path = fmt.Sprintf("%s%c", path, os.PathSeparator)
			_, err = w.Create(path)
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	//make sure path is relative to ZIP and not absolute (otherwise it'll make a bunch of nested folders)
	path := "./"
	err = filepath.Walk(path, walker)
	if err != nil {
		panic(err)
	}

	//close the ZIP file to make sure the uploader isn't uploading something broken
	w.Close()

	uploader = newUploader()
	upload()

	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)

}
