# MongoDB 模組範例專案

這是一個基於 Wails v2 框架的 MongoDB 操作範例專案，展示了如何在前端 Vue.js 應用中直接呼叫後端 Golang 函式來操作 MongoDB 資料庫。專案採用現代化的架構設計，包含資料庫遷移、環境變數管理、跨平台支援等企業級功能。

## 🚀 功能特色

- ✅ **完整的 CRUD 操作**：創建、讀取、更新、刪除用戶資料
- ✅ **分頁查詢功能**：高效處理大量數據
- ✅ **搜尋功能**：支援按姓名和電子郵件搜尋（使用正則表達式）
- ✅ **響應式前端介面**：現代化的 Vue 3 UI
- ✅ **模組化設計**：清晰的代碼結構和職責分離
- ✅ **資料庫遷移**：使用自定義 Go 遷移系統管理資料庫版本和索引
- ✅ **環境變數管理**：支援 .env 配置檔案
- ✅ **錯誤處理**：完整的錯誤處理和日誌記錄
- ✅ **跨平台支援**：支援 Windows、macOS、Linux
- ✅ **Singleton 模式**：確保資料庫實例的唯一性和線程安全
- ✅ **索引管理**：自動創建唯一索引和查詢索引

## 📁 專案結構

```
db-mongodb/
├── app.go                  # 主要應用邏輯和 API 方法
├── database.go             # MongoDB 資料庫操作模組
├── env_loader.go           # 環境變數載入模組
├── logger.go               # 日誌記錄模組
├── main.go                 # 應用程式入口點
├── go.mod                  # Go 模組依賴
├── go.sum                  # Go 依賴鎖定檔案
├── wails.json              # Wails 配置檔案
├── .env                    # 環境變數配置（需自行創建）
├── docker-compose.yml      # Docker Compose 配置（MongoDB）
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

提供完整的 MongoDB 資料庫操作功能：

- **Singleton 模式**：使用 `sync.Once` 確保資料庫實例的唯一性
- **連接管理**：自動處理資料庫連接的開啟和關閉
- **資料庫遷移**：使用自定義 Go 遷移系統，管理集合創建和索引
- **版本控制**：在 `migrations` 集合中記錄遷移版本
- **CRUD 操作**：完整的用戶增刪改查功能
- **分頁查詢**：支援分頁和搜尋功能
- **錯誤處理**：完整的錯誤處理和日誌記錄
- **索引管理**：自動創建唯一索引（email）和查詢索引（name, created_at）

**主要方法：**
```go
- Initialize()                                    // 初始化資料庫和執行遷移
- Connect()                                       // 連接到 MongoDB
- Disconnect()                                    // 關閉資料庫連接
- InsertUser(name, email, age)                   // 插入用戶
- GetAllUsers()                                  // 獲取所有用戶
- GetUserByID(id)                                // 根據 ID 獲取用戶（使用 ObjectID）
- UpdateUser(id, name, email, age)               // 更新用戶
- DeleteUser(id)                                 // 刪除用戶
- SearchUsers(keyword, page, pageSize)           // 搜尋用戶（支援分頁）
```

### 2. App 模組 (`app.go`)

提供前端可呼叫的 API 方法：

```go
- CreateUser(name, email, age)                   // 創建用戶
- GetAllUsers()                                  // 獲取所有用戶
- GetUser(id string)                             // 根據 ID 獲取用戶（ID 為字串）
- UpdateUser(id string, name, email, age)        // 更新用戶
- DeleteUser(id string)                          // 刪除用戶
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
- **MongoDB**：NoSQL 文件資料庫
- **mongo-driver**：MongoDB 官方 Go 驅動
- **joho/godotenv**：環境變數管理

### 前端技術

- **Vue 3**：響應式前端框架（Composition API）
- **Vite**：快速建置工具
- **JavaScript ES6+**：現代 JavaScript 語法

## 📦 安裝和運行

### 前置需求

- Go 1.23 或更高版本
- Node.js 16+ 和 npm
- MongoDB 4.4+ 資料庫
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### 1. 克隆專案

```bash
git clone <repository-url>
cd db-mongodb
```

### 2. 設定環境變數

創建 `.env` 檔案，然後填入您的資料庫連接資訊：

```env
# MongoDB Database Connection Settings
DB_HOST=localhost
DB_PORT=27017
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=mydb
DB_AUTH_SOURCE=admin
```

**環境變數說明：**
- `DB_HOST`：MongoDB 主機位址（預設：localhost）
- `DB_PORT`：MongoDB 連接埠（預設：27017）
- `DB_USER`：MongoDB 使用者名稱（可選，無認證時可留空）
- `DB_PASSWORD`：MongoDB 密碼（可選，無認證時可留空）
- `DB_NAME`：資料庫名稱（預設：mydb）
- `DB_AUTH_SOURCE`：認證資料庫（預設：admin）

### 3. 啟動 MongoDB 資料庫

**選項 A：使用 Docker Compose（推薦）**

```bash
docker-compose up -d
```

這會啟動一個帶有認證的 MongoDB 容器：
- 使用者名稱：`admin`
- 密碼：`admin`
- 資料庫：`mydb`

**選項 B：使用本地 MongoDB**

確保 MongoDB 服務正在運行。如果使用認證，請確保用戶有適當的權限。

**選項 C：使用 MongoDB Atlas（雲端）**

如果使用 MongoDB Atlas，請使用連接字串格式：
```env
DB_HOST=cluster0.xxxxx.mongodb.net
DB_PORT=27017
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=mydb
DB_AUTH_SOURCE=admin
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
2. 連接到 MongoDB 資料庫
3. 執行資料庫遷移（創建集合和索引）
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
   - **電子郵件**（必填，必須唯一）
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
3. 支援按姓名和電子郵件搜尋（不區分大小寫）

### 5. 分頁瀏覽

當用戶數量超過每頁顯示數量（預設 10 筆）時，系統會自動顯示分頁按鈕。

## 🔧 資料庫遷移

### MongoDB 遷移系統

與 SQL 資料庫不同，MongoDB 使用自定義的 Go 遷移系統來管理：

- **集合創建**：自動創建所需的集合
- **索引管理**：自動創建和更新索引
- **版本控制**：在 `migrations` 集合中記錄遷移版本

### 遷移檔案結構

遷移邏輯位於 `database.go` 中的 `runMigrations()` 方法：

```go
migrations := []func(context.Context) error{
    d.migration001_CreateUsersCollection, // Version 1
    // 未來可以在這裡添加更多遷移
    // d.migration002_AddNewIndex,        // Version 2
}
```

### 添加新的遷移

1. 在 `database.go` 中添加新的遷移函數：

```go
// migration002_AddNewIndex 添加新的索引
func (d *Database) migration002_AddNewIndex(ctx context.Context) error {
    collection := d.DB.Collection("users")
    
    indexModel := mongo.IndexModel{
        Keys: bson.D{
            {Key: "age", Value: 1},
        },
    }
    
    _, err := collection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return fmt.Errorf("failed to create index: %w", err)
    }
    
    WriteAppLog("Created age index on users collection", true)
    return nil
}
```

2. 在 `migrations` 切片中添加新遷移：

```go
migrations := []func(context.Context) error{
    d.migration001_CreateUsersCollection,
    d.migration002_AddNewIndex, // 新增
}
```

3. 重新啟動應用，遷移將自動執行

### 遷移狀態管理

- 應用啟動時會自動檢查並執行待處理的遷移
- 遷移版本記錄在 `migrations` 集合中
- 所有遷移操作都會記錄在日誌中
- 已執行的遷移不會重複執行

### 當前遷移內容

**Migration 001 - CreateUsersCollection：**
- 創建 `users` 集合（如果不存在）
- 創建唯一索引：`email`（確保電子郵件唯一性）
- 創建排序索引：`created_at`（用於按時間排序）
- 創建搜尋索引：`name`（用於姓名搜尋）

## 🔍 環境變數說明

| 變數名稱 | 說明 | 預設值 | 必填 |
|---------|------|--------|------|
| `DB_HOST` | 資料庫主機位址 | localhost | 否 |
| `DB_PORT` | 資料庫連接埠 | 27017 | 否 |
| `DB_USER` | 資料庫使用者名稱 | - | 否（無認證時可留空）|
| `DB_PASSWORD` | 資料庫密碼 | - | 否（無認證時可留空）|
| `DB_NAME` | 資料庫名稱 | mydb | 否 |
| `DB_AUTH_SOURCE` | 認證資料庫 | admin | 否 |

## 🚨 常見問題

### 1. 無法連接到資料庫

**問題：** `failed to ping MongoDB: connection refused`

**解決方案：**
- 確認 MongoDB 服務正在運行
- 檢查 `.env` 檔案中的連接資訊是否正確
- 確認防火牆設定允許連接到 MongoDB 埠（27017）
- 如果使用 Docker，確認容器正在運行：`docker ps`

### 2. 認證失敗

**問題：** `authentication failed`

**解決方案：**
- 確認 `DB_USER` 和 `DB_PASSWORD` 正確
- 確認 `DB_AUTH_SOURCE` 設定正確（通常是 `admin`）
- 如果使用 Docker Compose，預設用戶名和密碼都是 `admin`

### 3. 電子郵件重複錯誤

**問題：** `email already exists`

**解決方案：**
- 這是正常行為，系統會防止重複的電子郵件
- 如果需要更新用戶，請使用編輯功能

### 4. .env 檔案無法載入

**問題：** `.env file not found`

**解決方案：**
- 確認 `.env` 檔案存在於專案根目錄
- 檢查檔案編碼是否為 UTF-8（無 BOM）
- 應用會嘗試使用系統環境變數作為備選

### 5. 遷移失敗

**問題：** `failed to run migrations`

**解決方案：**
- 檢查 MongoDB 連接是否正常
- 確認用戶有創建集合和索引的權限
- 查看 `app.log` 檔案獲取詳細錯誤訊息

## 🎯 擴展建議

### 1. 添加更多集合

在 `database.go` 中添加新的遷移來創建其他集合：

```go
func (d *Database) migration002_CreateProductsCollection(ctx context.Context) error {
    collection := d.DB.Collection("products")
    
    indexes := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "name", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "price", Value: 1}},
        },
    }
    
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

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

### 5. 使用 MongoDB 聚合管道

利用 MongoDB 的聚合功能進行複雜查詢：

```go
func (d *Database) GetUsersByAgeRange(minAge, maxAge int) ([]map[string]interface{}, error) {
    pipeline := mongo.Pipeline{
        {{"$match", bson.D{
            {"age", bson.D{{"$gte", minAge}, {"$lte", maxAge}}},
        }}},
        {{"$sort", bson.D{{"created_at", -1}}}},
    }
    
    cursor, err := d.DB.Collection("users").Aggregate(ctx, pipeline)
    // ...
}
```

## 📝 注意事項

1. **資料庫連接**：使用 MongoDB 官方驅動，連接會在應用生命週期內保持
2. **並發安全**：Database 模組使用 Singleton 模式和 `sync.Once`，確保線程安全
3. **錯誤處理**：所有資料庫操作都包含完整的錯誤處理和日誌記錄
4. **密碼安全**：日誌中的密碼會被遮罩處理
5. **遷移管理**：請勿手動修改 `migrations` 集合
6. **跨平台**：專案支援 Windows、macOS、Linux 平台
7. **ObjectID**：MongoDB 使用 ObjectID 作為文檔 ID，在 API 中轉換為字串
8. **索引優化**：根據查詢模式適當添加索引以提升性能

## 🔄 從 PostgreSQL 遷移到 MongoDB

如果您是從 PostgreSQL 版本遷移過來的，請注意以下差異：

1. **ID 類型**：從 `int` 改為 `string`（ObjectID 的十六進制字串）
2. **查詢語法**：從 SQL 改為 MongoDB 查詢語法
3. **遷移系統**：從 SQL 檔案改為 Go 函數
4. **連接字串**：使用 MongoDB URI 格式
5. **索引管理**：在 Go 代碼中定義，而非 SQL

## 📄 授權

此專案僅供學習和參考使用。

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📧 聯絡方式

如有問題或建議，請聯絡專案維護者。
