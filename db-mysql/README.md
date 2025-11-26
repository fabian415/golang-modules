# MySQL 模組範例專案

這是一個基於 Wails v2 框架的 MySQL 操作範例專案，展示了如何在前端 Vue.js 應用中直接呼叫後端 Golang 函式來操作 MySQL 資料庫。專案採用現代化的架構設計，包含資料庫遷移、環境變數管理、跨平台支援等企業級功能。

## 🚀 功能特色

- ✅ **完整的 CRUD 操作**：創建、讀取、更新、刪除用戶資料
- ✅ **分頁查詢功能**：高效處理大量數據
- ✅ **搜尋功能**：支援按姓名和電子郵件搜尋
- ✅ **響應式前端介面**：現代化的 Vue 3 UI
- ✅ **模組化設計**：清晰的代碼結構和職責分離
- ✅ **資料庫遷移**：使用 golang-migrate 管理資料庫版本
- ✅ **環境變數管理**：支援 .env 配置檔案
- ✅ **錯誤處理**：完整的錯誤處理和日誌記錄
- ✅ **跨平台支援**：支援 Windows、macOS、Linux
- ✅ **Singleton 模式**：確保資料庫實例的唯一性和線程安全

## 📁 專案結構

```
golang_module_postgres/
├── app.go                  # 主要應用邏輯和 API 方法
├── database.go             # PostgreSQL 資料庫操作模組
├── env_loader.go           # 環境變數載入模組
├── logger.go               # 日誌記錄模組
├── main.go                 # 應用程式入口點
├── go.mod                  # Go 模組依賴
├── go.sum                  # Go 依賴鎖定檔案
├── wails.json              # Wails 配置檔案
├── .env                    # 環境變數配置（需自行創建）
├── .env.sample             # 環境變數範例檔案
├── docker-compose.yml      # Docker Compose 配置
├── _assets/                # 資源檔案
│   └── db/
│       └── migration/      # 資料庫遷移檔案
│           ├── 1_create_users_table.up.sql
│           └── 1_create_users_table.down.sql
├── frontend/               # Vue.js 前端
│   ├── src/
│   │   ├── App.vue         # 主要 Vue 組件
│   │   ├── main.js         # Vue 應用入口
│   │   └── style.css       # 全域樣式
│   ├── package.json        # 前端依賴
│   └── vite.config.js      # Vite 配置
└── build/                  # 建置輸出目錄
```

## 🏗️ 核心模組說明

### 1. Database 模組 (`database.go`)

提供完整的 PostgreSQL 資料庫操作功能：

- **Singleton 模式**：使用 `sync.Once` 確保資料庫實例的唯一性
- **連接管理**：自動處理資料庫連接的開啟和關閉
- **資料庫遷移**：自動執行 SQL 遷移檔案，支援版本管理
- **Dirty 狀態修復**：自動檢測並修復異常的遷移狀態
- **CRUD 操作**：完整的用戶增刪改查功能
- **分頁查詢**：支援分頁和搜尋功能
- **錯誤處理**：完整的錯誤處理和日誌記錄

**主要方法：**
```go
- Initialize()                                    // 初始化資料庫和執行遷移
- OpenDB() (*sql.DB, error)                      // 開啟資料庫連接
- InsertUser(name, email, age)                   // 插入用戶
- GetAllUsers()                                  // 獲取所有用戶
- GetUserByID(id)                                // 根據 ID 獲取用戶
- UpdateUser(id, name, email, age)               // 更新用戶
- DeleteUser(id)                                 // 刪除用戶
- SearchUsers(keyword, page, pageSize)           // 搜尋用戶（支援分頁）
```

### 2. App 模組 (`app.go`)

提供前端可呼叫的 API 方法：

```go
- CreateUser(name, email, age)                   // 創建用戶
- GetAllUsers()                                  // 獲取所有用戶
- GetUser(id)                                    // 根據 ID 獲取用戶
- UpdateUser(id, name, email, age)               // 更新用戶
- DeleteUser(id)                                 // 刪除用戶
- SearchUsers(keyword, page, pageSize)           // 搜尋用戶（支援分頁）
```

### 3. 環境變數載入模組 (`env_loader.go`)

負責載入 `.env` 檔案中的環境變數：

- 使用 `joho/godotenv` 套件
- 支援自訂 `.env` 檔案路徑
- 提供檔案存在性檢查
- 詳細的錯誤訊息

### 4. 前端介面 (`frontend/src/App.vue`)

提供完整的用戶管理介面：

- **表單操作**：新增和編輯用戶
- **數據展示**：表格形式展示用戶列表
- **搜尋功能**：即時搜尋用戶
- **分頁功能**：支援分頁瀏覽
- **響應式設計**：適配不同螢幕尺寸
- **狀態管理**：使用 Vue 3 Composition API

## 🛠️ 技術架構

### 後端技術

- **Go 1.23+**：主要程式語言
- **Wails v2.9.2**：跨平台桌面應用框架
- **MySQL 8.0+**：企業級關聯式資料庫
- **go-sql-driver/mysql**：MySQL 驅動
- **golang-migrate/migrate**：資料庫遷移工具
- **joho/godotenv**：環境變數管理

### 前端技術

- **Vue 3**：響應式前端框架（Composition API）
- **Vite**：快速建置工具
- **JavaScript ES6+**：現代 JavaScript 語法

## 📦 安裝和運行

### 前置需求

- Go 1.23 或更高版本
- Node.js 16+ 和 npm
- MySQL 8.0+ 資料庫
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### 1. 克隆專案

```bash
git clone <repository-url>
cd golang_module_postgres
```

### 2. 設定環境變數

複製 `.env.sample` 並重命名為 `.env`，然後填入您的資料庫連接資訊：

```bash
cp .env.sample .env
```

編輯 `.env` 檔案：

```env
# MySQL Database Connection Settings
DB_HOST=localhost
DB_PORT=3306
DB_USER=user
DB_PASSWORD=password
DB_NAME=mydb
```

### 3. 啟動 MySQL 資料庫

**選項 A：使用 Docker Compose（推薦）**

```bash
docker-compose up -d
```

**選項 B：使用本地 MySQL**

確保 MySQL 服務正在運行，並創建對應的資料庫：

```sql
CREATE DATABASE mydb;
```

### 4. 安裝依賴

```bash
# 安裝 Go 依賴
go mod tidy

# 安裝前端依賴
cd frontend
npm install
cd ..
```

### 5. 開發模式運行

```bash
wails dev
```

應用程式將自動：
1. 載入 `.env` 環境變數
2. 連接到 MySQL 資料庫
3. 執行資料庫遷移
4. 啟動前端開發伺服器
5. 開啟桌面應用視窗

### 6. 建置生產版本

```bash
wails build
```

建置完成後，可執行檔案將位於 `build/bin/` 目錄下。

## 📖 使用說明

### 1. 創建用戶

1. 在表單中填寫用戶資訊：
   - **姓名**（必填）
   - **電子郵件**（必填）
   - **年齡**（可選）
2. 點擊「創建」按鈕即可新增用戶

### 2. 編輯用戶

1. 點擊用戶列表中的「編輯」按鈕
2. 表單會自動填入該用戶的資訊
3. 修改後點擊「更新」即可

### 3. 刪除用戶

1. 點擊用戶列表中的「刪除」按鈕
2. 確認後即可刪除用戶

### 4. 搜尋用戶

1. 在搜尋框中輸入關鍵字
2. 點擊「搜尋」按鈕
3. 支援按姓名和電子郵件搜尋

### 5. 分頁瀏覽

當用戶數量超過每頁顯示數量（預設 10 筆）時，系統會自動顯示分頁按鈕。

## 🔧 資料庫遷移

### 遷移檔案結構

遷移檔案位於 `_assets/db/migration/` 目錄：

- `1_create_users_table.up.sql` - 創建 users 表
- `1_create_users_table.down.sql` - 刪除 users 表

### 添加新的遷移

1. 在 `_assets/db/migration/` 目錄下創建新的遷移檔案：
   ```
   2_add_new_feature.up.sql
   2_add_new_feature.down.sql
   ```

2. 編寫 SQL 語句

3. 重新啟動應用，遷移將自動執行

### 遷移狀態管理

- 應用啟動時會自動檢查並執行待處理的遷移
- 如果遷移處於 dirty 狀態，系統會自動嘗試修復
- 所有遷移操作都會記錄在日誌中

## 🔍 環境變數說明

| 變數名稱 | 說明 | 預設值 | 必填 |
|---------|------|--------|------|
| `DB_HOST` | 資料庫主機位址 | localhost | 否 |
| `DB_PORT` | 資料庫連接埠 | 3306 | 否 |
| `DB_USER` | 資料庫使用者名稱 | root | 否 |
| `DB_PASSWORD` | 資料庫密碼 | - | 是 |
| `DB_NAME` | 資料庫名稱 | mydb | 否 |

## 🚨 常見問題

### 1. 無法連接到資料庫

**問題：** `failed to ping database: connection refused`

**解決方案：**
- 確認 MySQL 服務正在運行
- 檢查 `.env` 檔案中的連接資訊是否正確
- 確認防火牆設定允許連接到 MySQL 埠（3306）
- 確認 MySQL 用戶有足夠的權限訪問資料庫

### 2. 遷移失敗

**問題：** `failed to run migrations: Dirty database version`

**解決方案：**
- 應用會自動嘗試修復 dirty 狀態
- 如果自動修復失敗，可以手動連接資料庫並檢查 `schema_migrations` 表

### 3. .env 檔案無法載入

**問題：** `.env file not found`

**解決方案：**
- 確認 `.env` 檔案存在於專案根目錄
- 檢查檔案編碼是否為 UTF-8（無 BOM）
- 參考 `.env.sample` 創建 `.env` 檔案

## 🎯 擴展建議

### 1. 添加更多資料表

在 `_assets/db/migration/` 目錄下創建新的遷移檔案：

```sql
-- 2_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

在 `database.go` 中添加對應的操作方法。

### 2. 添加資料驗證

在 `app.go` 中添加更嚴格的資料驗證：

```go
import "regexp"

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func (a *App) CreateUser(name, email string, age int) (string, error) {
    if !isValidEmail(email) {
        return "", fmt.Errorf("invalid email format")
    }
    
    if age < 0 || age > 150 {
        return "", fmt.Errorf("invalid age range")
    }
    
    // 繼續創建用戶邏輯...
}
```

### 3. 添加身份驗證

考慮添加用戶身份驗證和授權機制：
- JWT Token 驗證
- 密碼加密（bcrypt）
- 角色權限管理

### 4. 添加單元測試

為核心功能添加單元測試：

```go
// database_test.go
func TestInsertUser(t *testing.T) {
    db := GetDBInstance()
    err := db.InsertUser("Test User", "test@example.com", 25)
    if err != nil {
        t.Errorf("InsertUser failed: %v", err)
    }
}
```

## 📝 注意事項

1. **資料庫連接**：每次資料庫操作都會開啟新連接並在操作完成後關閉，避免資源洩漏
2. **並發安全**：Database 模組使用 Singleton 模式和 `sync.Once`，確保線程安全
3. **錯誤處理**：所有資料庫操作都包含完整的錯誤處理和日誌記錄
4. **密碼安全**：日誌中的密碼會被遮罩處理
5. **遷移管理**：請勿手動修改 `schema_migrations` 表
6. **跨平台**：專案支援 Windows、macOS、Linux 平台

## 📄 授權

此專案僅供學習和參考使用。

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📧 聯絡方式

如有問題或建議，請聯絡專案維護者。