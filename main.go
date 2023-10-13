package main

import (
	"context"
	"fmt"
	"go-storage/azureutil"
	"log"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func main() {

	// ストレージアカウントの接続文字列を作成
	connStr := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net", config.StorageAccountName, config.StorageAccountKey)

	// Blobサービスクライアントを作成
	containerURL, err := azureutil.GetContainerURL(connStr, config.ContainerName)
	if err != nil {
		log.Fatal(err)
	}

	// 画像ファイルを開く
	file, err := os.Open("Snap/Sample1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Blobに画像をアップロード
	blobURL := containerURL.NewBlockBlobURL("Sample1.jpg")
	_, err = azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpeg",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image uploaded successfully!")
}
