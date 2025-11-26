package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvFile 負責載入 .env 檔案中的環境變數
// 這回答了你的問題：os.Getenv() 本身不讀取 .env，它只讀取系統環境變數。
// 我們需要像 godotenv 這樣的套件將 .env 的內容載入到系統環境變數中。
func LoadEnvFile() error {
	// 預設路徑為當前目錄下的 .env
	envPath := ".env"

	// 如果你想指定特定路徑，可以在這裡修改，例如：
	// envPath = "./config/.env" 
	// envPath = "C:\\Users\\fabian.chung\\my_configs\\.env"

	// 檢查檔案是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// 如果檔案不存在，回傳錯誤（或者你可以選擇忽略，視需求而定）
		return fmt.Errorf(".env file not found at: %s", envPath)
	}

	// 使用 godotenv.Load 載入指定路徑的 .env 檔案
	// 這會將檔案中的變數注入到 process 的環境變數中
	// 之後你就可以使用 os.Getenv("KEY") 來取得這些值
	err := godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("error loading .env file from %s: %w", envPath, err)
	}

	return nil
}
