package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
	"log"
	"os"
	"path"
)

func uploadFileWeb3Storage(filename string) error {
	token := os.Getenv("WEB3_STORAGE_TOKEN")
	if token == "" {
		return fmt.Errorf("specify WEB3_STORAGE_TOKEN to upload your files to web3.storage")
	}

	client, errClient := w3s.NewClient(w3s.WithToken(token))
	if errClient != nil {
		return errClient
	}

	// Open the file for reading
	file, errFile := os.Open(filename)
	if errFile != nil {
		return errFile
	}

	basename := path.Base(filename)
	cidFile, errUpload := client.Put(context.Background(), file)
	if errUpload != nil {
		return errUpload
	}

	cidV0, _ := cid.Parse(cidFile)

	gatewayURL := fmt.Sprintf("https://w3s.link/ipfs/%s/%s", cidV0, basename)
	log.Printf("[INFO] Uploaded %s -> '%s'", filename, gatewayURL)

	return nil
}
