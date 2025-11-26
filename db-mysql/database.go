package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/go-sql-driver/mysql"
)

var SQL_MIGRATION_PATH string

// Database 結構體封裝所有資料庫操作
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// 使用 sync.Once 來確保 Singleton 實例只被創建一次
var once sync.Once
var dbInstance *Database

// GetDBInstance 獲取 Singleton 實例
func GetDBInstance() *Database {
	once.Do(func() {
		// 讀取環境變數並記錄
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "3306")
		user := getEnv("DB_USER", "root")
		password := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "mydb")
		
		// 記錄讀取到的配置（密碼只顯示長度）
		passwordMask := "***"
		if password != "" {
			passwordMask = fmt.Sprintf("*** (%d chars)", len(password))
		}
		WriteAppLog(fmt.Sprintf("Database config - Host: %s, Port: %s, User: %s, Password: %s, DB: %s", 
			host, port, user, passwordMask, dbName), true)
		
		dbInstance = &Database{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbName,
		}
		dbInstance.Initialize()
	})
	return dbInstance
}

// getEnv 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


// Initialize 初始化資料庫
func (d *Database) Initialize() {
	SQL_MIGRATION_PATH = "_assets/db/migration"

	if err := d.runMigrations(); err != nil {
		WriteAppLog(fmt.Sprintf("Failed to run migrations: %v", err), true)
	}
}

// OpenDB 開啟資料庫連接
func (d *Database) OpenDB() (*sql.DB, error) {
	// MySQL DSN 格式: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// 測試連接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

// region <-- Migration 相關函式-->
// 執行資料庫 migration
func (d *Database) runMigrations() error {
	WriteAppLog("Initializing database migration...", true)
	WriteAppLog(fmt.Sprintf("DB Host: %s, Port: %s, DB: %s", d.Host, d.Port, d.DBName), true)
	WriteAppLog("Migration Path: "+SQL_MIGRATION_PATH, true)

	// Convert relative path to absolute path to avoid issues with migrate
	absMigrationPath, err := filepath.Abs(SQL_MIGRATION_PATH)
	if err != nil {
		return fmt.Errorf("failed to get absolute migration path: %w", err)
	}
	
	// Ensure path uses forward slashes for file:// URL on Windows
	absMigrationPath = filepath.ToSlash(absMigrationPath)

	// 構建 MySQL 連接字串（golang-migrate 格式）
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
		d.User, d.Password, d.Host, d.Port, d.DBName)

	m, err := migrate.New(
		"file://"+absMigrationPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// 檢查當前版本
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		WriteAppLog("Database is empty, will run all migrations", true)
	} else {
		WriteAppLog(fmt.Sprintf("Current database version: %d, dirty: %t", version, dirty), true)

		// 如果資料庫是 dirty 狀態，嘗試修復
		if dirty {
			WriteAppLog("Database is dirty, attempting to fix...", true)
			if err := d.fixDirtyDatabase(m, version); err != nil {
				return fmt.Errorf("failed to fix dirty database: %w", err)
			}
		}
	}

	maxVer, err := maxMigrationVersion(SQL_MIGRATION_PATH)
	if err != nil {
		return fmt.Errorf("failed to scan migration dir: %w", err)
	}
	WriteAppLog(fmt.Sprintf("Max migration file %d", maxVer), true)

	if int(version) > maxVer {
		// 預設不 downgrade
		WriteAppLog(fmt.Sprintf("DB version %d > max migration file %d — migration files missing", version, maxVer), true)
	} else {
		// 檢查是否有可用的 migration
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				WriteAppLog("Database is already up to date", true)
				return nil
			}
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	// 再次檢查版本確認更新成功
	newVersion, newDirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	WriteAppLog(fmt.Sprintf("Database successfully migrated to version: %d, dirty: %t", newVersion, newDirty), true)

	return nil
}

func maxMigrationVersion(migPath string) (int, error) {
	entries, err := os.ReadDir(migPath)
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)
	max := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := re.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		v, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if v > max {
			max = v
		}
	}
	return max, nil
}

// 修復 dirty 狀態的資料庫
func (d *Database) fixDirtyDatabase(m *migrate.Migrate, currentVersion uint) error {
	// 方法1: 強制設定版本為當前版本（清除 dirty 標記）
	if err := m.Force(int(currentVersion)); err != nil {
		return fmt.Errorf("failed to force version %d: %w", currentVersion, err)
	}

	WriteAppLog(fmt.Sprintf("Successfully forced database version to %d", currentVersion), true)
	return nil
}
// endregion

// InsertUser 插入用戶資料
func (d *Database) InsertUser(name, email string, age int) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	insertSQL := `INSERT INTO users (name, email, age) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, name, email, age)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// GetAllUsers 獲取所有用戶
func (d *Database) GetAllUsers() ([]map[string]interface{}, error) {
	db, err := d.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		user := make(map[string]interface{})
		for i, column := range columns {
			user[column] = values[i]
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID 根據 ID 獲取用戶
func (d *Database) GetUserByID(id int) (map[string]interface{}, error) {
	db, err := d.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	user := make(map[string]interface{})
	for i, column := range columns {
		user[column] = values[i]
	}

	return user, nil
}

// UpdateUser 更新用戶資料
func (d *Database) UpdateUser(id int, name, email string, age int) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	updateSQL := `UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?`
	result, err := db.Exec(updateSQL, name, email, age, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser 刪除用戶
func (d *Database) DeleteUser(id int) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM users WHERE id = ?`
	result, err := db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers 搜尋用戶（支援分頁）
func (d *Database) SearchUsers(keyword string, page, pageSize int) ([]map[string]interface{}, int, error) {
	db, err := d.OpenDB()
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	// 計算總數
	var count int
	countSQL := `SELECT COUNT(*) FROM users WHERE name LIKE ? OR email LIKE ?`
	searchPattern := "%" + keyword + "%"
	err = db.QueryRow(countSQL, searchPattern, searchPattern).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 查詢分頁資料
	offset := (page - 1) * pageSize
	querySQL := `SELECT * FROM users WHERE name LIKE ? OR email LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(querySQL, searchPattern, searchPattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get columns: %w", err)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		user := make(map[string]interface{})
		for i, column := range columns {
			user[column] = values[i]
		}
		users = append(users, user)
	}

	return users, count, nil
}
