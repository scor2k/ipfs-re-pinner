# ipfs-re-pinner
Tool to re-pin your CIDs from one IPFS node to another

## how to build

```shell
git clone git@github.com:scor2k/ipfs-re-pinner.git
cd ipfs-re-pinner/
go mod tidy
go build -o ipfs-re-pinner
```

## how to use
```shell
./ipfs-re-pinner re-pin --old https://old-ipfs-server.io:5001 --new https://old-ipfs-server.io:5001 --cid QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK
2023/02/08 21:11:23 [INFO] File QmbMQkxgyCVDZkxcm1Sx5a7BhZ2ZmvQvXrfmshoYUSirbK.png (image/png) was successfully re-pinned to the new IPFS node
```
