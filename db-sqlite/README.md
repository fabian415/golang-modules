# SQLite 模組範例專案

這是一個基於 Wails v2 框架的 SQLite 操作範例專案，展示了如何在前端 Vue.js 應用中直接呼叫後端 Golang 函式來操作 SQLite 資料庫。

## 功能特色

- ✅ 完整的 CRUD 操作（創建、讀取、更新、刪除）
- ✅ 分頁查詢功能
- ✅ 搜尋功能
- ✅ 響應式前端介面
- ✅ 模組化設計
- ✅ 錯誤處理

## 專案結構

```
example/
├── app.go              # 主要應用邏輯和 API 方法
├── database.go         # SQLite 資料庫操作模組
├── main.go            # 應用程式入口點
├── go.mod             # Go 模組依賴
├── frontend/          # Vue.js 前端
│   ├── src/
│   │   ├── App.vue    # 主要 Vue 組件
│   │   └── main.js    # Vue 應用入口
│   └── package.json   # 前端依賴
└── README.md          # 專案說明
```

## 核心模組說明

### 1. Database 模組 (`database.go`)

提供完整的 SQLite 資料庫操作功能：

- **Singleton 模式**：確保資料庫實例的唯一性
- **連接管理**：自動處理資料庫連接的開啟和關閉
- **CRUD 操作**：支援用戶的增刪改查
- **分頁查詢**：支援分頁和搜尋功能
- **錯誤處理**：完整的錯誤處理機制

### 2. App 模組 (`app.go`)

提供前端可呼叫的 API 方法：

- `CreateUser(name, email, age)` - 創建用戶
- `GetAllUsers()` - 獲取所有用戶
- `GetUser(id)` - 根據 ID 獲取用戶
- `UpdateUser(id, name, email, age)` - 更新用戶
- `DeleteUser(id)` - 刪除用戶
- `SearchUsers(keyword, page, pageSize)` - 搜尋用戶（支援分頁）

### 3. 前端介面 (`App.vue`)

提供完整的用戶管理介面：

- **表單操作**：新增和編輯用戶
- **數據展示**：表格形式展示用戶列表
- **搜尋功能**：支援按姓名和電子郵件搜尋
- **分頁功能**：支援分頁瀏覽
- **響應式設計**：適配不同螢幕尺寸

## 安裝和運行

### 1. 安裝依賴

```bash
# 安裝 Go 依賴
go mod tidy

# 安裝前端依賴
cd frontend
npm install
cd ..
```

### 2. 開發模式運行

```bash
wails dev
```

### 3. 建置生產版本

```bash
wails build
```

## 使用說明

### 1. 創建用戶

在表單中填寫用戶資訊：
- 姓名（必填）
- 電子郵件（必填）
- 年齡（可選）

點擊「創建」按鈕即可新增用戶。

### 2. 編輯用戶

點擊用戶列表中的「編輯」按鈕，表單會自動填入該用戶的資訊，修改後點擊「更新」即可。

### 3. 刪除用戶

點擊用戶列表中的「刪除」按鈕，確認後即可刪除用戶。

### 4. 搜尋用戶

在搜尋框中輸入關鍵字，點擊「搜尋」按鈕即可搜尋用戶。支援按姓名和電子郵件搜尋。

### 5. 分頁瀏覽

當用戶數量較多時，系統會自動顯示分頁按鈕，點擊即可切換頁面。

## 技術架構

### 後端技術

- **Go 1.21+**：主要程式語言
- **Wails v2**：桌面應用框架
- **SQLite3**：輕量級資料庫
- **mattn/go-sqlite3**：SQLite 驅動
- **golang-migrate**：資料庫遷移工具

### 前端技術

- **Vue 3**：響應式前端框架
- **Vite**：快速建置工具
- **JavaScript ES6+**：現代 JavaScript 語法

## 擴展建議

### 1. 添加更多資料表

可以在 `database.go` 中添加更多資料表的操作，例如：

```go
// 創建產品表
func (d *Database) CreateProductTable() error {
    // 實現產品表創建邏輯
}

// 產品 CRUD 操作
func (d *Database) InsertProduct(name, description string, price float64) error {
    // 實現產品插入邏輯
}
```

### 2. 添加資料驗證

可以在 `app.go` 中添加更嚴格的資料驗證：

```go
func (a *App) CreateUser(name, email string, age int) (string, error) {
    // 驗證電子郵件格式
    if !isValidEmail(email) {
        return "", fmt.Errorf("invalid email format")
    }
    
    // 驗證年齡範圍
    if age < 0 || age > 150 {
        return "", fmt.Errorf("invalid age range")
    }
    
    // 繼續創建用戶邏輯...
}
```

### 3. 添加日誌功能

可以添加日誌記錄功能來追蹤操作：

```go
import "log"

func (a *App) CreateUser(name, email string, age int) (string, error) {
    log.Printf("Creating user: %s, %s, %d", name, email, age)
    
    // 創建用戶邏輯...
    
    log.Printf("User created successfully: %s", email)
    return "User created successfully", nil
}
```

## 注意事項

1. **資料庫檔案**：SQLite 資料庫檔案會創建在專案根目錄下的 `example.db`
2. **並發安全**：Database 模組使用 Singleton 模式，確保資料庫操作的線程安全
3. **錯誤處理**：所有資料庫操作都包含完整的錯誤處理
4. **資源管理**：資料庫連接會自動關閉，避免資源洩漏

## 授權

此專案僅供學習和參考使用。