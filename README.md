# async-blob-copy

Simple go command to copy blobs between storage accounts

## usage

```golang
$ ./async-blob-copy \
-sa <source storage account name> \
-sk <source storage account key> \
-sc <source container name> \
-da <destination storage account name> \
-dk <destination storage account key> \
-dc <destination container name> \
-b  <blob name>
```
