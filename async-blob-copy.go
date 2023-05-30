package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"

	"github.com/schollz/progressbar/v3"
)

var (
	srcAccountName   string
	srcAccountKey    string
	srcContainerName string
	dstAccountName   string
	dstContainerName string
	dstAccountKey    string
	blobName         string
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	flag.StringVar(&srcAccountName, "sa", srcAccountName, "source account name")
	flag.StringVar(&srcAccountKey, "sk", srcAccountKey, "source account key")
	flag.StringVar(&srcContainerName, "sc", srcContainerName, "source container name")
	flag.StringVar(&dstAccountName, "da", dstAccountName, "destination account name")
	flag.StringVar(&dstAccountKey, "dk", dstAccountKey, "destination account key")
	flag.StringVar(&dstContainerName, "dc", dstContainerName, "destination container name")
	flag.StringVar(&blobName, "b", blobName, "blob name")
	flag.Parse()

	startCopyWithAccountKey(srcAccountName, srcContainerName, dstAccountName, dstContainerName, blobName)
}

func getSASToken(accountName string, containerName string, blobName string, credential *azblob.SharedKeyCredential) (url string, token string, err error) {
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC().Add(-48 * time.Hour),
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   to.Ptr(sas.ContainerPermissions{Read: true, Create: true, Write: true}).String(),
		ContainerName: containerName,
		BlobName:      blobName,
	}.SignWithSharedKey(credential)
	handleError(err)

	token = fmt.Sprintf("%s", sasQueryParams.Encode())
	url = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	return url, token, nil
}

func startCopyWithAccountKey(srcAccountName string, srcContainerName string, dstAccountName string, dstContainerName string, blobName string) {

	ctx := context.Background()
	dstUrl, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", dstAccountName, dstContainerName, blobName))
	handleError(err)

	srcSharedKeyCredential, err := azblob.NewSharedKeyCredential(srcAccountName, srcAccountKey)
	handleError(err)

	dstSharedKeyCredential, err := azblob.NewSharedKeyCredential(dstAccountName, dstAccountKey)
	handleError(err)

	uri, token, err := getSASToken(srcAccountName, srcContainerName, blobName, srcSharedKeyCredential)
	sasTokenURI := fmt.Sprintf("%s?%s", uri, token)
	handleError(err)

	src, _ := url.Parse(sasTokenURI)
	blockBlobClient, err := blockblob.NewClientWithSharedKeyCredential(dstUrl.String(), dstSharedKeyCredential, nil)
	handleError(err)

	copyOperation, err := blockBlobClient.StartCopyFromURL(ctx, src.String(), nil)
	handleError(err)

	copyID := copyOperation.CopyID
	copyStatus := copyOperation.CopyStatus

	fmt.Printf("copy: %s --> %s\nid:%s\n", uri, blockBlobClient.URL(), *copyID)

	bar := progressbar.NewOptions(100, progressbar.OptionSetPredictTime(true))

	for *copyStatus == *copyOperation.CopyStatus {
		
		getMetadata, err := blockBlobClient.GetProperties(ctx, nil)
		handleError(err)

		copyStatus = getMetadata.CopyStatus

		copyCompleted, err := strconv.ParseFloat(strings.Split(*getMetadata.CopyProgress, "/")[0], 64)
		handleError(err)

		copyTotal, err := strconv.ParseFloat(strings.Split(*getMetadata.CopyProgress, "/")[1], 64)
		handleError(err)

		percentComplete := (copyCompleted / copyTotal) * 100
		// fmt.Printf("progress: %0.2f %%\n", percentComplete)
		bar.Set(int(percentComplete))
		time.Sleep(500 * time.Millisecond)
	}
	bar.Finish()
	fmt.Print("\n")
}
