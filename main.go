package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	SFTP_USER     = "sftp user here"
	SFTP_PASSWORD = "sftp password here"
	SFTP_HOST     = "server ip/host:22" // Dont forget the port (22)! i.e. 127.0.0.1:22
	SFTP_FOLDER   = "folder name"       // Folder name where the CSV file will be uploaded to, i.e. <SSH HOME>/my-folder
)

func main() {
	fmt.Println("Generating a report...")
	fileName := generateReport()

	fmt.Println("Uploading file to the server...")
	uploadFileSFTP(fileName)
}

// Generate report
func generateReport() string {
	// Generate filename
	currentTime := time.Now().Local()
	currentDate := currentTime.Format("2006-01-02")

	var fileName = currentDate + "_my_file.csv"

	// Create/open/Chmod the file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)

	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	csvWriter := csv.NewWriter(file)

	// Heading of the CSV
	heading := [][]string{[]string{"Country", "City"}}

	csvWriter.WriteAll(heading)

	// Rows of the CSV -his could be a for loop, i.e. Read data from DB - for each row in DB - add to CSV
	csvWriter.WriteAll([][]string{[]string{"Ireland", "Dublin"}})
	csvWriter.WriteAll([][]string{[]string{"Lithuania", "Vilnius"}})

	csvWriter.Flush()

	fmt.Println("Report has been generating asuccessfully!")

	return fileName
}

// Upload file through SFTP
func uploadFileSFTP(fileName string) {
	// Set up Sftp connection
	config := &ssh.ClientConfig{
		User:            SFTP_USER,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(SFTP_PASSWORD),
		},
	}

	config.SetDefaults()

	sshConn, err := ssh.Dial("tcp", SFTP_HOST, config)

	if err != nil {
		log.Fatal("Connection failed", err)
	}

	defer sshConn.Close()

	client, err := sftp.NewClient(sshConn)

	if err != nil {
		log.Fatal("Failed to create sftp client: ", err)
	}

	dstFile, err := client.Create(SFTP_FOLDER + "/" + fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer dstFile.Close()

	srcFile, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.Copy(dstFile, io.LimitReader(srcFile, 3e9))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully uploaded " + fileName + " to the server.")
	log.Printf("Wrote %d bytes\n", bytes)
}
