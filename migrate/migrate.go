package main

import (
	"fmt"
	"go-storage/azureutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 環境構築用storageアカウントにContainer作成する
func main() {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load() //.env読み込み
		if err != nil {
			log.Fatalln(err) //強制終了
		}
	}
	client := azureutil.CreateContainer()
	defer fmt.Println("Successfully Migrated")
	// azureutil.UploadBlob(client)
	// azureutil.UploadSnap(client)
	// azureutil.DownloadSnap(client)
	azureutil.GetContainerList(client)
}
