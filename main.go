package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
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
			{
				Name:  "upload-web3",
				Usage: "upload file to web3.storage service",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "file",
						Required: true,
						Usage:    "path to the file to upload",
					},
				},
				Action: func(c *cli.Context) error {
					uploadError := uploadFileWeb3Storage(c.Path("file"))
					if uploadError != nil {
						log.Fatalf("[ERROR] %+v", uploadError)
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
