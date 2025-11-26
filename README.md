# Golang Modules 使用說明

本專案包含四個基於 Golang + Wails 的模組專案：**custom-script**、**db-sqlite**、**db-mysql**、**db-postgres**。每個專案皆可獨立運作，也可作為大型跨平台應用的模組範本。

## 📦 專案概述

這四個模組展示了使用 Wails v2 框架開發跨平台桌面應用的不同場景：

- **custom-script**: 執行自訂bash/bat腳本並顯示建置進度的應用
- **db-sqlite**: 使用 SQLite 資料庫的 CRUD 操作與資料庫遷移 (database migration) 管理範例
- **db-mysql**: 使用 MySQL 資料庫的 CRUD 操作與資料庫遷移 (database migration) 管理範例
- **db-postgres**: 使用 PostgreSQL 資料庫的 CRUD 操作與資料庫遷移 (database migration) 管理範例

## 🚀 快速開始

### 前置需求

- **Go 1.23+** 或更高版本
- **Node.js 16+** 和 npm
- **Wails CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- **MySQL 8.0+**（僅 db-mysql 模組需要）
- **PostgreSQL 12+**（僅 db-postgres 模組需要）

### 安裝 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 選擇模組

每個模組都是獨立的專案，可以單獨使用：

```bash
# 進入想要使用的模組目錄
cd custom-script    # 或 db-sqlite 或 db-mysql 或 db-postgres

# 安裝前端依賴
cd frontend
npm install
cd ..

# 開發模式運行
wails dev

# 建置生產版本
wails build
```
## 🎯 模組詳細說明

### 1. custom-script

**功能特色：**
- ✅ 執行自訂建置腳本（支援 Windows 和 Unix 系統）
- ✅ 即時顯示建置進度
- ✅ 即時日誌輸出
- ✅ 支援取消建置操作
- ✅ 自動捲動日誌
- ✅ 建置狀態顯示（成功/失敗/進行中）

**使用場景：**
- 自動化建置流程
- CI/CD 工具整合
- 批次處理任務監控

**詳細文檔：** 請參考 [custom-script/README.md](custom-script/README.md)

### 2. db-sqlite

**功能特色：**
- ✅ 完整的 CRUD 操作（創建、讀取、更新、刪除）
- ✅ 分頁查詢功能
- ✅ 搜尋功能
- ✅ 響應式前端介面
- ✅ 模組化設計
- ✅ 錯誤處理

**使用場景：**
- 輕量級資料庫應用
- 單機桌面應用
- 原型開發和測試

**詳細文檔：** 請參考 [db-sqlite/README.md](db-sqlite/README.md)

### 3. db-postgres

**功能特色：**
- ✅ 完整的 CRUD 操作
- ✅ 分頁查詢和搜尋功能
- ✅ 資料庫遷移管理（golang-migrate）
- ✅ 環境變數管理（.env 支援）
- ✅ Docker Compose 支援
- ✅ 完整的錯誤處理和日誌記錄
- ✅ 跨平台支援（Windows、macOS、Linux）

**使用場景：**
- 企業級桌面應用
- 需要多用戶協作的應用
- 需要複雜資料庫操作的應用

**詳細文檔：** 請參考 [db-postgres/README.md](db-postgres/README.md)

### 4. db-mysql

**功能特色：**
- ✅ 完整的 CRUD 操作
- ✅ 分頁查詢和搜尋功能
- ✅ 資料庫遷移管理（golang-migrate）
- ✅ 環境變數管理（.env 支援）
- ✅ Docker Compose 支援
- ✅ 完整的錯誤處理和日誌記錄
- ✅ 跨平台支援（Windows、macOS、Linux）
- ✅ Singleton 模式：確保資料庫實例的唯一性和線程安全

**使用場景：**
- 企業級桌面應用
- 需要多用戶協作的應用
- 需要複雜資料庫操作的應用
- 使用 MySQL 作為後端資料庫的應用

**詳細文檔：** 請參考 [db-mysql/README.md](db-mysql/README.md)

## 🛠️ 技術架構

### 後端技術

- **Go 1.23+**: 主要程式語言
- **Wails v2**: 跨平台桌面應用框架
- **SQLite3**: 輕量級資料庫（db-sqlite）
- **MySQL 8.0+**: 企業級關聯式資料庫（db-mysql）
- **PostgreSQL**: 企業級關聯式資料庫（db-postgres）
- **golang-migrate**: 資料庫遷移工具（db-mysql、db-postgres）
- **go-sql-driver/mysql**: MySQL 驅動（db-mysql）
- **joho/godotenv**: 環境變數管理（db-mysql、db-postgres）

### 前端技術

- **Vue 3**: 響應式前端框架（Composition API）
- **Vite**: 快速建置工具
- **JavaScript ES6+**: 現代 JavaScript 語法

## 📝 開發指南

### 添加新模組

1. 在根目錄創建新的模組目錄
2. 使用 `wails init` 初始化新專案
3. 參考現有模組的結構和模式
4. 更新本 README.md 文件

### 模組間共享代碼

如果需要共享代碼，可以考慮：

1. 創建共享的 Go 模組
2. 使用 Go Workspace（`go.work`）
3. 將共享代碼提取到獨立的 package

### 建置和部署

```bash
# 開發模式（自動重新編譯和熱重載）
wails dev

# 建置生產版本
wails build

# 建置特定平台
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

## 🔧 環境配置

### db-mysql 環境變數

如果使用 db-mysql 模組，需要配置 `.env` 檔案：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=user
DB_PASSWORD=your_password
DB_NAME=mydb
```

### db-postgres 環境變數

如果使用 db-postgres 模組，需要配置 `.env` 檔案：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mydb
DB_SSLMODE=disable
```

### 使用 Docker Compose

**db-mysql:**
```bash
cd db-mysql
docker-compose up -d
```

**db-postgres:**
```bash
cd db-postgres
docker-compose up -d
```

## 📚 學習資源

- [Wails 官方文檔](https://wails.io/docs/)
- [Vue 3 官方文檔](https://vuejs.org/)
- [Go 官方文檔](https://go.dev/doc/)

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 授權

此專案僅供學習和參考使用。

## 📧 聯絡方式

如有問題或建議，請聯絡專案維護者。
