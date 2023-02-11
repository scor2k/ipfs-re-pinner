package main

import (
	"bytes"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	ipfsShell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := &cli.App{
		Name:     "ipfs-re-pinner",
		HelpName: "re-pin your CIDs from one IPFS node to another",
		Usage:    "Just use me!",
		Flags:    []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "re-pin",
				Usage: "specify old and new servers + CID",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "timeout",
						Required: false,
						Value:    int64(10),
						Usage:    "HTTP Request timeout (GET only, seconds)",
					},
					&cli.StringFlag{
						Name:     "old",
						Required: true,
						Usage:    "for example: https://old-ipfs-server.io:5001",
					},
					&cli.StringFlag{
						Name:     "new",
						Required: true,
						Usage:    "for example: https://new-ipfs-server.io:5001",
					},
					&cli.StringFlag{
						Name:     "cid",
						Required: true,
						Usage:    "CID hash",
					}},
				Action: func(c *cli.Context) error {
					errPin := rePinCID(c.String("old"), c.String("new"), c.String("cid"), c.Int64("timeout"))
					if errPin != nil {
						log.Printf("[ERROR] %+v", errPin)
					}
					return nil
				},
			},
			{
				Name:  "download",
				Usage: "download CID from IPFS site",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ipfs",
						Required: true,
						Usage:    "for example: https://old-ipfs-server.io:5001",
					},
					&cli.StringFlag{
						Name:     "dir",
						Required: true,
						Usage:    "directory to save CID",
					},
					&cli.StringFlag{
						Name:     "cid",
						Required: true,
						Usage:    "CID hash",
					},
					&cli.Int64Flag{
						Name:     "timeout",
						Required: false,
						Value:    int64(10),
						Usage:    "HTTP Request timeout (seconds)",
					},
				},
				Action: func(c *cli.Context) error {
					data, _, fileExt, errGet := getIPFS(c.String("ipfs"), c.String("cid"), c.Int64("timeout"))
					if errGet != nil {
						log.Fatalf("[ERROR] %+v", errGet)
					}

					fileName := filepath.Join(c.String("dir"), fmt.Sprintf("%s%s", c.String("cid"), fileExt))

					errSaveFile := os.WriteFile(fileName, data, 0644)
					if errSaveFile != nil {
						log.Fatalf("[ERROR] %+v", errGet)
					}
					log.Printf("[INFO] File %s was successfully saved", fileName)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// rePinCID - download file / upload file to the new IPFS / pin file on the new IPFS
func rePinCID(oldIPFSHost string, newIPFSHost string, CID string, getTimeout int64) error {
	data, fileType, fileExtension, errGet := getIPFS(oldIPFSHost, CID, getTimeout)
	if errGet != nil {
		return errGet
	}
	newCID, errUpload := uploadIPFS(newIPFSHost, data)
	if errUpload != nil {
		return errUpload
	}

	if newCID != CID {
		return fmt.Errorf("CID was changed %s -> %s", CID, newCID)
	}

	log.Printf("[INFO] File %s%s (%s) was successfully re-pinned to the new IPFS node", CID, fileExtension, fileType)
	return nil
}

// getIPFS - download file & return bytes, type, extension
func getIPFS(ipfsHost string, CID string, timeout int64) ([]byte, string, string, error) {
	shell := ipfsShell.NewShell(ipfsHost)
	shell.SetTimeout(time.Duration(timeout) * time.Second)

	urlCID := fmt.Sprintf("/ipfs/%s", CID)
	file, fileErr := shell.Cat(urlCID)

	if fileErr != nil {
		return nil, "", "", fileErr
	}

	body, bodyError := io.ReadAll(file)
	if bodyError != nil {
		return nil, "", "", bodyError
	}

	mimeType := mimetype.Detect(body)
	return body, mimeType.String(), mimeType.Extension(), nil
}

// uploadIPFS - upload bytes to IPFS, pin it and return the CID (should be the same)
func uploadIPFS(ipfsHost string, data []byte) (string, error) {
	shell := ipfsShell.NewShell(ipfsHost)
	shell.SetTimeout(300 * time.Second)
	ipfsReader := bytes.NewReader(data)

	ipfsCID, errIpfs := shell.Add(ipfsReader)
	if errIpfs != nil {
		return "", errIpfs
	}

	pinError := shell.Pin(ipfsCID)
	if pinError != nil {
		return "", pinError
	}
	return ipfsCID, nil
}
