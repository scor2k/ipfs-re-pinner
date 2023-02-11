# ipfs-re-pinner
Tool to re-pin your CIDs from one IPFS node to another

Also, it can help you to:
 - download file from IPFS
 - upload file to IPFS via https://web3.storage service (provide your token via ENV)

## how to build

```shell
git clone git@github.com:scor2k/ipfs-re-pinner.git
cd ipfs-re-pinner/
go mod tidy
go build -o ipfs-re-pinner
```

## Examples

### how to re-pin CIDs
```shell
./ipfs-re-pinner re-pin --old https://old-ipfs-server.io:5001 --new https://old-ipfs-server.io:5001 --cid QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK
2023/02/08 21:11:23 [INFO] File QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK.png (image/png) was successfully re-pinned to the new IPFS node
```

### how to download CIDs
```shell
./ipfs-re-pinner download --ipfs http://dummy-ipfs-server:5001 --dir backup --cid QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK
2023/02/09 21:15:03 [INFO] File backup/QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK.png was successfully saved
```

### how to set up request timeout
```shell
./ipfs-re-pinner download --timeout 5 --ipfs http://dummy-ipfs-server:5001 --dir backup --cid QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK
                          ^^^^^^^^^^^^
                        timeout 5 seconds
```

### hot to upload files to the web3.storage service
```shell
export WEB3_STORAGE_TOKEN=aaaaaaaa
./ipfs-re-pinner upload-web3 --file ~/tmp/ipfs/Qma2uX5KD2W9nwYjVdi6CQJabb5a245o5Qn8JBCE97eAMm.png
```
