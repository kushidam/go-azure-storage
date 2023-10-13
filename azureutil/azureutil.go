package azureutil

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ConnectionCiient() *azblob.Client {
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("'AZURE_STORAGE_CONNECTION_STRING' not found")
	}
	// 接続文字列でクライアントを作成する
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	handleError(err)
	return client
}

func CreateContainer() *azblob.Client {
	ctx := context.Background()
	client := ConnectionCiient()
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	resp, err := client.CreateContainer(ctx, containerName, &azblob.CreateContainerOptions{})
	if err != nil {
		if bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
			// コンテナが既に存在する
			fmt.Println("container already exists")
		} else {
			handleError(err)
		}
	}
	fmt.Println(resp)
	return client
}

func UploadBlob(client *azblob.Client) {
	ctx := context.Background()
	data := []byte("\nHello, world! This is a blob.\n")
	blobName := "sample-blob"
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	// Upload to data to blob storage
	fmt.Printf("Uploading a blob named %s\n", blobName)
	_, err := client.UploadBuffer(ctx, containerName, blobName, data, &azblob.UploadBufferOptions{})
	handleError(err)
}

func UploadSnap(client *azblob.Client) {
	ctx := context.Background()
	snapFile, err := os.Open("./snap/Sample1.jpg")
	handleError(err)
	snapFileName := "Sample1.jpg"
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	// Upload to data to blob storage
	fmt.Printf("Uploading a blob named %s\n", snapFileName)
	_, err = client.UploadFile(ctx, containerName, snapFileName, snapFile, &azblob.UploadBufferOptions{})
	handleError(err)
}

func DownloadSnap(client *azblob.Client) {
	ctx := context.Background()
	snapFileName := "Sample1.jpg"
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	file, err := os.Create("./download/Sample1.jpg")
	handleError(err)
	defer file.Close()
	fmt.Printf("Downloading %s\n", snapFileName)
	_, err = client.DownloadFile(ctx, containerName, snapFileName, file, &azblob.DownloadFileOptions{})
	handleError(err)

}

func GetContainerList(client *azblob.Client) {
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	fmt.Println("Listing the blobs in the container:")
	// ページャーを作成して、すべてのページを反復処理する
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)

		for _, blob := range resp.Segment.BlobItems {
			// 各ブロブの名前を出力する
			fmt.Println(*blob.Name)
		}
	}
}
