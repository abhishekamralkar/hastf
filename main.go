package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {

	tfVersion := flag.String("v", "", "terraform version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *tfVersion == "" {
		fmt.Println("Error: Please provide a version to download.")
		flag.Usage()
		os.Exit(1)
	}

	if *tfVersion != "" {
		backupDir := "/tmp/terraform_backups"
		if _, err := exec.LookPath("terraform"); err == nil {
			fmt.Println("Terraform already installed. Creating backup...")

			if _, err := os.Stat(backupDir); os.IsNotExist(err) {
				os.Mkdir(backupDir, os.ModePerm)
			}

			backupFile := fmt.Sprintf("%s/terraform", backupDir)
			fmt.Println(backupFile)
			err := copyFile(backupFile, "/usr/local/bin/terraform")
			if err != nil {
				fmt.Println("Error creating backup:", err)
				os.Exit(1)
			}
			fmt.Println("Backup created at:", backupFile)
		}
		err := downloadFile(*tfVersion, "terraform.zip")
		if err != nil {
			fmt.Println("Error downloading Terraform:", err)
			os.Exit(1)
		}

		err = exec.Command("unzip", "terraform.zip").Run()
		if err != nil {
			fmt.Println("Error extracting Terraform:", err)
			os.Exit(1)
		}

		err = exec.Command("mv", "terraform", "/usr/local/bin/").Run()
		if err != nil {
			fmt.Println("Error moving Terraform:", err)
			os.Exit(1)
		}
		fmt.Printf("Terraform %s installed successfully.\n", *tfVersion)
	}
}

// downloadFile downloads a file from the specified URL and saves it locally
func downloadFile(tfVersion, filename string) error {
	tfURL := "https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip"
	url := fmt.Sprintf(tfURL, tfVersion, tfVersion, runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// copyFile copies a file from source to destination
func copyFile(dest, src string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
