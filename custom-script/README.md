# Custom Scripts Builder

這是一個完整的 Golang + Wails 專案，用於執行自訂腳本並顯示建置進度。

## 專案結構

```
custom-scripts/
├── app.go                    # 後端主應用邏輯
├── main.go                   # 應用入口
├── process_unix.go           # Unix 系統進程處理
├── process_windows.go        # Windows 系統進程處理
├── go.mod                    # Go 模組依賴
├── wails.json                # Wails 設定檔
├── scripts/                  # 腳本目錄
│   ├── build_script.bat      # Windows 建置腳本
│   └── build_script.sh       # Unix 建置腳本
└── frontend/                 # 前端目錄
    ├── src/
    │   ├── components/
    │   │   └── SampleBuildProgress.vue  # 建置進度元件
    │   ├── router/
    │   │   └── index.js      # 路由設定
    │   ├── App.vue           # 主應用元件
    │   ├── main.js           # 前端入口
    │   └── style.css         # 全域樣式
    ├── index.html            # HTML 範本
    ├── package.json          # 前端依賴
    └── vite.config.js        # Vite 設定
```

## 功能特性

- 執行自訂建置腳本（支援 Windows 和 Unix 系統）
- 即時顯示建置進度
- 即時日誌輸出
- 支援取消建置操作
- 自動捲動日誌
- 建置狀態顯示（成功/失敗/進行中）

## 開發環境設定

### 前置要求

- Go 1.23 或更高版本
- Node.js 和 npm
- Wails CLI

### 安裝步驟

1. 安裝 Wails CLI（如果尚未安裝）：
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

2. 安裝前端依賴：
```bash
cd frontend
npm install
cd ..
```

3. 安裝 Go 依賴：
```bash
go mod download
```

## 執行專案

### 開發模式

```bash
wails dev
```

這將啟動開發伺服器，自動重新編譯和熱重載。

### 建置生產版本

```bash
wails build
```

## 使用說明

1. 啟動應用後，點擊 "Start Build" 按鈕開始建置
2. 建置過程中可以查看即時進度和日誌
3. 可以隨時點擊 "Cancel Build" 取消建置
4. 日誌區域支援自動捲動，也可以手動捲動查看歷史日誌

## 腳本格式

建置腳本需要輸出特定格式的進度資訊：

- 進度百分比：`BSP Build Progress: XX%`
- 成功標記：`[SUCCESS]`
- 錯誤標記：`[ERROR]`
- 資訊標記：`[INFO]`

## 注意事項

- 腳本檔案需要放在 `scripts/` 目錄下
- Windows 使用 `build_script.bat`，Unix 系統使用 `build_script.sh`
- 建置日誌會保存在專案根目錄的 `sample_build.log` 檔案中

