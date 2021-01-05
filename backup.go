package main

import (
	"log"
	"time"
	"flag"
	"os"
	"fmt"

	"github.com/jlaffaye/ftp"
)

func uploadFile(c *ftp.ServerConn, file string, debug bool) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Print("Starting the upload - ", file)
	}
	err = c.Append("/backup" + fmt.Sprintf("%s", time.Now().Format("20060102150405")) + ".tar.xz", f)
	if debug {
		log.Print("The file is uploaded")
	}
}

func deleteRotateDays(maxDays int, c *ftp.ServerConn, debug bool) {
	dt := time.Now()
	entries, err := c.List("/")
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Print("Full file list:")
	}
	for _, entry := range entries {
		if debug {
			log.Println(entry.Name, " - ", entry.Time, " - ", dt.Sub(entry.Time).Hours() / 24)
		}
		if int(dt.Sub(entry.Time).Hours() / 24) > maxDays {
			if debug {
				log.Println("To delete - ", entry.Name)
			}
			c.Delete(entry.Name)
		} 
	}
} 

func main() {
	var (
		debug bool
		rotateDays int
		rotateFiles int
		file string
		host string
		password string
		user string
		backupPath string
	)

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.IntVar(&rotateDays, "rotateDays", 10, "rotate days")
	flag.IntVar(&rotateFiles, "rotateFiles", 0, "rotate files")
	flag.StringVar(&file, "file", "", "file to upload")
	flag.StringVar(&host, "host", "", "FTP host ex. ftp.example.com:21")
	flag.StringVar(&user, "user", "", "FTP user")
	flag.StringVar(&password, "password", "", "FTP password")
	flag.StringVar(&backupPath, "backupPath", "/", "FTP backup path")
	flag.Parse()

	if host == "" || file == "" {
		log.Fatal("Please check args.")
	}

	c, err := ftp.Dial(host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}
	
	if rotateDays != 0 {
		deleteRotateDays(rotateDays, c, debug)
	}

	//if rotateFiles != 0 {
		// todo
	//}

	uploadFile(c, file, debug)

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}