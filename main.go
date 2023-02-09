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
					errPin := rePinCID(c.String("old"), c.String("new"), c.String("cid"))
					if errPin != nil {
						log.Printf("[ERROR] %+v", errPin)
					}
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
func rePinCID(oldIPFSHost string, newIPFSHost string, CID string) error {
	data, fileType, fileExtension, errGet := getIPFS(oldIPFSHost, CID)
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
func getIPFS(ipfsHost string, CID string) ([]byte, string, string, error) {
	shell := ipfsShell.NewShell(ipfsHost)
	shell.SetTimeout(300 * time.Second)

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
