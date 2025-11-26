# MongoDB Migration 新版本操作指南

文件目的：示範如何在 `db-mongodb` 專案中新增一個 **可實際執行** 的遷移版本，並用專案內的工具鏈（Go、Wails、Docker Compose、mongosh）完成驗證。

---

## Step 0 — 建立測試環境

1. 於專案根目錄準備 `.env`（若尚未存在），內容參考：
   ```env
   DB_HOST=localhost
   DB_PORT=27017
   DB_USER=admin
   DB_PASSWORD=admin
   DB_NAME=mydb
   DB_AUTH_SOURCE=admin
   ```
2. 啟動 MongoDB（使用專案附的 Compose）：
   ```powershell
   docker-compose up -d
   ```
3. 驗證資料庫可連線：
   ```powershell
   mongosh "mongodb://admin:admin@localhost:27017/mydb?authSource=admin"
   ```

---

## Step 1 — 選定版本與名稱

1. 打開 `database.go`。
2. 找到 `runMigrations()` 中的 `migrations` slice，確定目前最新版本（預設為 `migration001_CreateUsersCollection`）。
3. 將下一個版本命名為 `migration002_<Feature>`，例如：`migration002_AddSignupAuditIndex`。

---

## Step 2 — 實作新的 Migration 函式

以下以「在 `users` 集合加上 `signup_channel` 欄位索引」為例。

1. 在 `// region <-- Migration 相關函式-->` 底下新增函式：

```200:249:database.go
func (d *Database) migration002_AddSignupAuditIndex(ctx context.Context) error {
    collection := d.DB.Collection("users")

    index := mongo.IndexModel{
        Keys: bson.D{{Key: "signup_channel", Value: 1}},
        Options: options.Index().
            SetName("signup_channel_idx").
            SetBackground(true),
    }

    _, err := collection.Indexes().CreateOne(ctx, index)
    if err != nil {
        return fmt.Errorf("failed to create signup_channel index: %w", err)
    }

    WriteAppLog("Migration 002: ensured signup_channel index", true)
    return nil
}
```

2. 原則：
   - **冪等性**：使用 `CreateOne` / `CreateMany` 建索引時，若索引已存在不會出錯。
   - **錯誤訊息**：使用 `fmt.Errorf("context: %w", err)` 便於追蹤。
   - 若需更新資料，可使用 `UpdateMany`，並加上過濾條件避免覆寫新資料。

---

## Step 3 — 註冊新的 Migration

1. 在 `runMigrations()` 的 slice 中加入此函式，並保持版本順序：

```150:179:database.go
migrations := []func(context.Context) error{
    d.migration001_CreateUsersCollection, // Version 1
    d.migration002_AddSignupAuditIndex,   // Version 2 (新增)
}
```

2. 儲存檔案後，遷移系統會自動在套用成功時透過 `setVersion()` 更新 `migrations` 集合。

---

## Step 4 — 執行並監看 Migration

1. 在專案根目錄執行：
   ```powershell
   wails dev
   ```
   或若只驗證後端：
   ```powershell
   go run .
   ```
2. 終端輸出應包含以下訊息：
   ```
   Connecting to MongoDB...
   Running migration 2...
   Migration 2 completed successfully
   ```
3. 若看到錯誤訊息，依提示修正後重新執行命令。

---

## Step 5 — 驗證 Migration 成果

1. 使用 `mongosh` 驗證版本號與索引：
   ```powershell
   mongosh "mongodb://admin:admin@localhost:27017/mydb?authSource=admin" --eval "db.migrations.findOne()"
   mongosh "mongodb://admin:admin@localhost:27017/mydb?authSource=admin" --eval "db.users.getIndexes()"
   ```
2. 若遷移包含資料轉換，建議執行下列查詢確認：
   ```powershell
   mongosh "mongodb://admin:admin@localhost:27017/mydb?authSource=admin" --eval "db.users.find({signup_channel: {\$exists: true}}).limit(5)"
   ```

---

## Step 6 — 加入測試與提交

1. 執行後端測試確保功能正常：
   ```powershell
   go test ./...
   ```
2. 撰寫提交訊息：
   ```powershell
   git add database.go
   git commit -m "feat: add signup channel migration"
   ```
3. 在 PR 或文件中描述遷移內容、風險、必要的手動驗證步驟。

---

## 常見 FAQ

- **如何回滾？** 目前僅提供前進式 migration。若要回復，請撰寫新的遷移函式進行反向操作並遞增版本。
- **是否可批次新增多個版本？** 可以，但建議一次新增一個函式並逐步驗證，避免無法追蹤的錯誤。
- **如何避免重複執行？** `runMigrations()` 會讀取 `migrations` 集合中的 `version` 欄位，只執行大於目前版本的函式。

---

照上述步驟，即可以可重現、可驗證的方式為 `db-mongodb` 專案新增任何 MongoDB schema 變動。若在實作過程遇到新情境，請補充至文件或提交說明中，維持團隊共用的最佳實踐。*** End Patch

